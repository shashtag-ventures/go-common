package gcp

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/iam/apiv1/iampb"
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/api/iterator"
)

func (s *gcpService) DeployService(ctx context.Context, karadaProjectID string, imageName string, serviceName string) (string, error) {
	// Apply a timeout to prevent hanging the deployment process
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	parent := fmt.Sprintf("projects/%s/locations/%s", s.projectID, s.region)
	serviceID := fmt.Sprintf("projects/%s/locations/%s/services/%s", s.projectID, s.region, serviceName)

	// Check if service exists first?
	// CreateServiceRequest vs UpdateServiceRequest.
	// Or use purely CreateService and catch "AlreadyExists" then Update?
	// For simplicity in MVP, we might assume Create, or just try Create.
	// If it exists, we must Update.

	// Let's try to Get it first.
	exists := false
	_, err := s.runClient.GetService(ctx, &runpb.GetServiceRequest{Name: serviceID})
	if err == nil {
		exists = true
	}

	// Logic to capture URI after operation
	var serviceURI string

	if exists {
		// UPDATE
		req := &runpb.UpdateServiceRequest{
			Service: &runpb.Service{
				Name: serviceID,
				Template: &runpb.RevisionTemplate{
					Containers: []*runpb.Container{
						{
							Image: imageName,
							Ports: []*runpb.ContainerPort{
								{ContainerPort: 8080},
							},
							Resources: &runpb.ResourceRequirements{
								CpuIdle:          true,
								StartupCpuBoost: true,
							},
						},
					},
					Scaling: &runpb.RevisionScaling{
						MaxInstanceCount: 5,
					},
					SessionAffinity: true,
				},
			},
		}
		op, err := s.runClient.UpdateService(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to update service: %w", err)
		}
		resp, err := op.Wait(ctx) // This waits for deployment to finish!
		if err != nil {
			return "", fmt.Errorf("update service operation failed: %w", err)
		}
		serviceURI = resp.Uri

	} else {
		// CREATE
		req := &runpb.CreateServiceRequest{
			Parent:    parent,
			ServiceId: serviceName,
			Service: &runpb.Service{
				Template: &runpb.RevisionTemplate{
					Containers: []*runpb.Container{
						{
							Image: imageName,
							Ports: []*runpb.ContainerPort{
								{ContainerPort: 8080},
							},
							Resources: &runpb.ResourceRequirements{
								CpuIdle:          true,
								StartupCpuBoost: true,
							},
						},
					},
					Scaling: &runpb.RevisionScaling{
						MinInstanceCount: 0,
						MaxInstanceCount: 5, 
					},
					SessionAffinity: true, // Crucial for stateful MCP sessions
				},
				Ingress: runpb.IngressTraffic_INGRESS_TRAFFIC_ALL,
			},
		}

		op, err := s.runClient.CreateService(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to create service: %w", err)
		}

		resp, err := op.Wait(ctx)
		if err != nil {
			return "", fmt.Errorf("create service operation failed: %w", err)
		}
		serviceURI = resp.Uri
	}

	// Make Public (Allow unauthenticated invocations)
	// Apply this for both Create and Update (in case permissions were changed or lost)
	policy := &iampb.SetIamPolicyRequest{
		Resource: serviceID,
		Policy: &iampb.Policy{
			Bindings: []*iampb.Binding{
				{
					Role:    "roles/run.invoker",
					Members: []string{"allUsers"},
				},
			},
		},
	}
	if _, err := s.runClient.SetIamPolicy(ctx, policy); err != nil {
		// Log error but don't fail deployment completely? Or fail?
		// Better to fail or warning.
		return serviceURI, fmt.Errorf("failed to make service public: %w", err)
	}

	return serviceURI, nil
}

func (s *gcpService) ListServices(ctx context.Context) ([]string, error) {
	// Apply a timeout to prevent hanging the list process
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	parent := fmt.Sprintf("projects/%s/locations/%s", s.projectID, s.region)
	req := &runpb.ListServicesRequest{
		Parent: parent,
	}

	var services []string
	it := s.runClient.ListServices(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list services: %w", err)
		}
		services = append(services, resp.Name)
	}
	return services, nil
}

func (s *gcpService) DeleteService(ctx context.Context, serviceName string) error {
	// Apply a timeout to prevent hanging the delete process
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Service name format: projects/{project}/locations/{location}/services/{service}
	// The ListServices returns full names.
	// If input is just the short name, we construct full name.
	// But let's assume input is full name or handle both?
	// The implementation plan says "Use runClient.DeleteService".

	// Let's construct the full ID if it doesn't look like one, just to be safe,
	// OR assume caller passes full ID if they got it from ListServices.
	// But user might pass just "karada-xyz".
	// Let's rely on ListServices returning full names and we use those.

	req := &runpb.DeleteServiceRequest{
		Name: serviceName,
	}

	op, err := s.runClient.DeleteService(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete service %s: %w", serviceName, err)
	}

	if _, err := op.Wait(ctx); err != nil {
		return fmt.Errorf("delete service operation failed: %w", err)
	}

	return nil
}
