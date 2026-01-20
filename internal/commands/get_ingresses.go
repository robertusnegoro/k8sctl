package commands

import (
	"context"
	"fmt"

	"github.com/robertusnegoro/k8ctl/internal/errors"
	"github.com/robertusnegoro/k8ctl/internal/output"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getIngresses(client kubernetes.Interface, namespace, name, outputFormat string, showNamespace bool) error {
	var ingresses []interface{}

	if name != "" {
		ing, err := client.NetworkingV1().Ingresses(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "ingress", name, namespace)
		}
		ingresses = []interface{}{ing}
	} else {
		ingList, err := client.NetworkingV1().Ingresses(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return errors.HandleKubernetesError(err, "ingresses", "", namespace)
		}
		for i := range ingList.Items {
			ingresses = append(ingresses, &ingList.Items[i])
		}
	}

	if outputFormat == OutputFormatJSON {
		return output.FormatJSON(ingresses)
	}
	if outputFormat == OutputFormatYAML {
		return output.FormatYAML(ingresses)
	}

	headers := []string{"NAME", "CLASS", "HOSTS", "ADDRESS", "PORTS", "AGE"}
	if showNamespace {
		headers = []string{"NAMESPACE", "NAME", "CLASS", "HOSTS", "ADDRESS", "PORTS", "AGE"}
	}
	table := output.NewTable(headers)

	for _, ing := range ingresses {
		i, ok := ing.(*networkingv1.Ingress)
		if !ok {
			continue
		}

		age := getAge(i.CreationTimestamp)

		// Get ingress class
		class := NoneValue
		if i.Spec.IngressClassName != nil {
			class = *i.Spec.IngressClassName
		} else if i.Annotations["kubernetes.io/ingress.class"] != "" {
			class = i.Annotations["kubernetes.io/ingress.class"]
		}

		// Get hosts
		hosts := []string{}
		for _, rule := range i.Spec.Rules {
			if rule.Host != "" {
				hosts = append(hosts, rule.Host)
			}
		}
		hostsStr := NoneValue
		if len(hosts) > 0 {
			hostsStr = hosts[0]
			if len(hosts) > 1 {
				hostsStr = fmt.Sprintf("%s +%d more", hostsStr, len(hosts)-1)
			}
		}

		// Get address
		address := NoneValue
		if len(i.Status.LoadBalancer.Ingress) > 0 {
			lb := i.Status.LoadBalancer.Ingress[0]
			if lb.IP != "" {
				address = lb.IP
			} else if lb.Hostname != "" {
				address = lb.Hostname
			}
		}

		// Get ports
		ports := []string{}
		for _, rule := range i.Spec.Rules {
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					if path.Backend.Service != nil {
						port := fmt.Sprintf("%d", path.Backend.Service.Port.Number)
						ports = append(ports, port)
					}
				}
			}
		}
		portsStr := NoneValue
		if len(ports) > 0 {
			portsStr = ports[0]
			if len(ports) > 1 {
				portsStr = fmt.Sprintf("%s +%d more", portsStr, len(ports)-1)
			}
		}

		row := []string{
			i.Name,
			class,
			hostsStr,
			address,
			portsStr,
			age,
		}
		if showNamespace {
			row = []string{
				i.Namespace,
				i.Name,
				class,
				hostsStr,
				address,
				portsStr,
				age,
			}
		}
		table.AddRow(row)
	}

	table.Render()
	return nil
}
