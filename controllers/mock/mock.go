//Package mock_test contains mock clients for testing
package mock_test

//go:generate go run github.com/golang/mock/mockgen -package mock_test -destination client.go sigs.k8s.io/controller-runtime/pkg/client Client
//go:generate go run github.com/golang/mock/mockgen -package mock_test -destination pomerium_ingress_reconciler.go github.com/pomerium/ingress-controller/pomerium IngressReconciler
//go:generate go run github.com/golang/mock/mockgen -package mock_test -destination pomerium_config_reconciler.go github.com/pomerium/ingress-controller/pomerium ConfigReconciler
