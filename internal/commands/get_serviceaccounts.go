package commands

import (
	"context"
	"fmt"

	"github.com/robertusnegoro/k8ctl/internal/errors"
	"github.com/robertusnegoro/k8ctl/internal/output"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getServiceAccounts(client kubernetes.Interface, namespace, name, outputFormat string, showNamespace bool) error {
	var serviceAccounts []interface{}

	if name != "" {
		sa, err := client.CoreV1().ServiceAccounts(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "serviceaccount", name, namespace)
		}
		serviceAccounts = []interface{}{sa}
	} else {
		saList, err := client.CoreV1().ServiceAccounts(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "serviceaccounts", "", namespace)
		}
		for i := range saList.Items {
			serviceAccounts = append(serviceAccounts, &saList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(serviceAccounts)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(serviceAccounts)
	}

	headers := []string{"NAME", "SECRETS", "AGE"}
	if showNamespace {
		headers = []string{"NAMESPACE", "NAME", "SECRETS", "AGE"}
	}
	table := output.NewTable(headers)

	for _, sa := range serviceAccounts {
		s, ok := sa.(*corev1.ServiceAccount)
		if !ok {
			continue
		}

		age := getAge(s.CreationTimestamp)
		secretsCount := fmt.Sprintf("%d", len(s.Secrets))

		row := []string{
			s.Name,
			secretsCount,
			age,
		}
		if showNamespace {
			row = []string{
				s.Namespace,
				s.Name,
				secretsCount,
				age,
			}
		}
		table.AddRow(row)
	}

	table.Render()
	return nil
}
