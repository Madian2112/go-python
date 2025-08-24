# Computación en la Nube con Go

## Introducción

Go se ha convertido en uno de los lenguajes preferidos para el desarrollo de aplicaciones y servicios en la nube debido a su eficiencia, simplicidad y excelente soporte para la concurrencia. En este documento, exploraremos cómo utilizar Go para desarrollar aplicaciones nativas de la nube, interactuar con proveedores de servicios en la nube y aplicar patrones de diseño específicos para entornos cloud.

## Desarrollo de Aplicaciones Nativas de la Nube

### Principios de Diseño Cloud-Native

1. **Microservicios**: Arquitectura basada en servicios pequeños, independientes y acoplados libremente.
2. **Contenedores**: Empaquetado de aplicaciones y sus dependencias para garantizar consistencia en diferentes entornos.
3. **Orquestación**: Gestión automatizada de contenedores (típicamente con Kubernetes).
4. **Infraestructura como Código (IaC)**: Definición y provisión de infraestructura mediante código.
5. **Observabilidad**: Monitoreo, logging y trazabilidad integrados.
6. **Resiliencia**: Diseño para fallos y recuperación automática.

### Contenedorización con Docker

Go es ideal para crear aplicaciones contenedorizadas debido a su capacidad para compilar binarios estáticos:

```go
// main.go - Aplicación web simple
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "¡Hola desde Go en un contenedor!")
	})

	log.Printf("Servidor iniciado en puerto %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
```

Dockerfile para la aplicación:

```dockerfile
# Etapa de compilación
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Etapa final
FROM alpine:3.18

WORKDIR /app

# Copiar binario desde la etapa de compilación
COPY --from=builder /app/app .

# Exponer puerto
EXPOSE 8080

# Ejecutar aplicación
CMD ["./app"]
```

Comandos para construir y ejecutar:

```bash
# Construir imagen
docker build -t mi-app-go .

# Ejecutar contenedor
docker run -p 8080:8080 mi-app-go
```

### Desarrollo para Kubernetes

Go es ampliamente utilizado para desarrollar aplicaciones y operadores para Kubernetes:

#### Ejemplo de Cliente Kubernetes en Go

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Configuración del cliente
	var config *rest.Config
	var err error

	// Intentar cargar configuración desde dentro del clúster
	config, err = rest.InClusterConfig()
	if err != nil {
		// Si falla, intentar cargar desde kubeconfig
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error al cargar configuración: %v", err)
		}
	}

	// Crear cliente
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error al crear cliente: %v", err)
	}

	// Listar pods en el namespace default
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error al listar pods: %v", err)
	}

	fmt.Printf("Encontrados %d pods en el namespace default\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Printf("Pod: %s\n", pod.Name)
	}
}
```

#### Ejemplo de Operador Kubernetes

Los operadores de Kubernetes extienden la funcionalidad del clúster para gestionar aplicaciones específicas:

```go
package main

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Definición simplificada de un operador
func main() {
	// Obtener configuración
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	// Crear manager
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		panic(err)
	}

	// Crear controlador
	c, err := controller.New("mi-controlador", mgr, controller.Options{
		Reconciler: &MiReconciler{client: mgr.GetClient()},
	})
	if err != nil {
		panic(err)
	}

	// Observar recursos
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		panic(err)
	}

	// Iniciar manager
	if err := mgr.Start(context.Background()); err != nil {
		panic(err)
	}
}

// Reconciler para procesar eventos
type MiReconciler struct {
	client client.Client
}

func (r *MiReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	// Lógica de reconciliación
	fmt.Printf("Reconciliando %s\n", req.NamespacedName)
	return reconcile.Result{}, nil
}
```

### Serverless con Go

Go es excelente para funciones serverless debido a su tiempo de inicio rápido y bajo consumo de memoria:

#### AWS Lambda con Go

```go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extraer parámetros
	name := request.QueryStringParameters["name"]
	if name == "" {
		name = "Mundo"
	}

	// Procesar solicitud
	message := fmt.Sprintf("¡Hola, %s!", name)

	// Devolver respuesta
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Body: message,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
```

#### Google Cloud Functions con Go

```go
package function

import (
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloWorld", helloWorld)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	// Extraer parámetros
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Mundo"
	}

	// Devolver respuesta
	fmt.Fprintf(w, "¡Hola, %s!", name)
}
```

## Interacción con Proveedores de Servicios en la Nube

### AWS SDK para Go

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Cargar configuración de AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Crear cliente S3
	client := s3.NewFromConfig(cfg)

	// Listar buckets
	resp, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("Error al listar buckets: %v", err)
	}

	fmt.Println("Buckets:")
	for _, bucket := range resp.Buckets {
		fmt.Printf("- %s, creado el: %s\n", *bucket.Name, bucket.CreationDate)
	}
}
```

### Google Cloud SDK para Go

```go
package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()

	// Crear cliente de Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error al crear cliente: %v", err)
	}
	defer client.Close()

	// Listar buckets
	it := client.Buckets(ctx, "mi-proyecto")
	fmt.Println("Buckets:")
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error al listar buckets: %v", err)
		}
		fmt.Printf("- %s, creado el: %s\n", bucketAttrs.Name, bucketAttrs.Created)
	}
}
```

### Azure SDK para Go

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func main() {
	// Obtener credenciales
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Error al obtener credenciales: %v", err)
	}

	// Crear cliente de Blob Storage
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	client, err := azblob.NewClient(serviceURL, cred, nil)
	if err != nil {
		log.Fatalf("Error al crear cliente: %v", err)
	}

	// Listar contenedores
	ctx := context.Background()
	pager := client.NewListContainersPager(nil)

	fmt.Println("Contenedores:")
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Error al listar contenedores: %v", err)
		}

		for _, container := range resp.ContainerItems {
			fmt.Printf("- %s\n", *container.Name)
		}
	}
}
```

## Patrones de Diseño para la Nube

### Patrón Circuit Breaker

Evita que una aplicación intente repetidamente una operación que probablemente fallará:

```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

// Función que simula una llamada a un servicio externo
func llamarServicio() (string, error) {
	// Simular fallo aleatorio
	// En un caso real, aquí llamaríamos a un servicio externo
	return "", errors.New("servicio no disponible")
}

func main() {
	// Configurar circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "mi-servicio",
		MaxRequests: 3,                  // Número máximo de solicitudes permitidas cuando el circuito está medio abierto
		Interval:    5 * time.Second,     // Intervalo de tiempo para reiniciar contador de fallos
		Timeout:     30 * time.Second,    // Tiempo que el circuito permanece abierto
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("Circuit breaker '%s' cambió de '%s' a '%s'\n", name, from, to)
		},
	})

	// Intentar llamadas al servicio
	for i := 0; i < 10; i++ {
		result, err := cb.Execute(func() (interface{}, error) {
			return llamarServicio()
		})

		if err != nil {
			fmt.Printf("Intento %d: Error - %v\n", i+1, err)
		} else {
			fmt.Printf("Intento %d: Éxito - %v\n", i+1, result)
		}

		time.Sleep(1 * time.Second)
	}
}
```

### Patrón Retry con Backoff Exponencial

Reintenta operaciones fallidas con intervalos crecientes:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// Función que simula una operación que puede fallar
func operacionInestable() error {
	// Simular fallo aleatorio
	if rand.Float64() < 0.7 {
		return errors.New("error temporal")
	}
	return nil
}

func main() {
	// Configurar backoff exponencial
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 100 * time.Millisecond
	expBackoff.MaxInterval = 10 * time.Second
	expBackoff.MaxElapsedTime = 1 * time.Minute

	// Función a ejecutar con retry
	operation := func() error {
		err := operacionInestable()
		if err != nil {
			fmt.Printf("Operación falló: %v, reintentando...\n", err)
			return err
		}
		fmt.Println("Operación exitosa")
		return nil
	}

	// Ejecutar con retry
	ctx := context.Background()
	err := backoff.Retry(operation, backoff.WithContext(expBackoff, ctx))
	if err != nil {
		fmt.Printf("Error después de múltiples intentos: %v\n", err)
	} else {
		fmt.Println("Operación completada con éxito")
	}
}
```

### Patrón Bulkhead

Aísla componentes para evitar que los fallos se propaguen:

```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// Servicio con limitación de concurrencia
type Servicio struct {
	nombre    string
	semaphore *semaphore.Weighted
}

func NuevoServicio(nombre string, maxConcurrencia int64) *Servicio {
	return &Servicio{
		nombre:    nombre,
		semaphore: semaphore.NewWeighted(maxConcurrencia),
	}
}

func (s *Servicio) Ejecutar(ctx context.Context) error {
	// Intentar adquirir un slot
	if err := s.semaphore.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("no se pudo adquirir semáforo: %w", err)
	}
	defer s.semaphore.Release(1)

	// Simular procesamiento
	duracion := time.Duration(100+rand.Intn(900)) * time.Millisecond
	fmt.Printf("[%s] Iniciando operación (duración: %v)\n", s.nombre, duracion)
	time.Sleep(duracion)

	// Simular error aleatorio
	if rand.Float64() < 0.2 {
		return fmt.Errorf("error en el servicio %s", s.nombre)
	}

	fmt.Printf("[%s] Operación completada con éxito\n", s.nombre)
	return nil
}

func main() {
	// Crear servicios con diferentes límites de concurrencia
	servicioA := NuevoServicio("ServicioA", 2) // Máximo 2 operaciones concurrentes
	servicioB := NuevoServicio("ServicioB", 5) // Máximo 5 operaciones concurrentes

	// Ejecutar múltiples operaciones
	var wg sync.WaitGroup
	ctx := context.Background()

	// Simular 10 solicitudes al ServicioA
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := servicioA.Ejecutar(ctx)
			if err != nil {
				fmt.Printf("[ServicioA-%d] Error: %v\n", id, err)
			}
		}(i)
	}

	// Simular 10 solicitudes al ServicioB
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := servicioB.Ejecutar(ctx)
			if err != nil {
				fmt.Printf("[ServicioB-%d] Error: %v\n", id, err)
			}
		}(i)
	}

	// Esperar a que todas las operaciones terminen
	wg.Wait()
	fmt.Println("Todas las operaciones han terminado")
}
```

## Infraestructura como Código (IaC) con Go

### Pulumi con Go

Pulumi permite definir infraestructura como código utilizando Go:

```go
package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Crear una VPC
		vpc, err := ec2.NewVpc(ctx, "mi-vpc", &ec2.VpcArgs{
			CidrBlock: pulumi.String("10.0.0.0/16"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("mi-vpc"),
			},
		})
		if err != nil {
			return err
		}

		// Crear subredes públicas
		subnet1, err := ec2.NewSubnet(ctx, "subnet-1", &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String("10.0.1.0/24"),
			AvailabilityZone: pulumi.String("us-west-2a"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("subnet-1"),
			},
		})
		if err != nil {
			return err
		}

		subnet2, err := ec2.NewSubnet(ctx, "subnet-2", &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String("10.0.2.0/24"),
			AvailabilityZone: pulumi.String("us-west-2b"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("subnet-2"),
			},
		})
		if err != nil {
			return err
		}

		// Crear grupo de seguridad
		sg, err := ec2.NewSecurityGroup(ctx, "web-sg", &ec2.SecurityGroupArgs{
			VpcId: vpc.ID(),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(80),
					ToPort:     pulumi.Int(80),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String("web-sg"),
			},
		})
		if err != nil {
			return err
		}

		// Crear instancias EC2
		server1, err := ec2.NewInstance(ctx, "server-1", &ec2.InstanceArgs{
			Ami:                      pulumi.String("ami-0c55b159cbfafe1f0"),
			InstanceType:             pulumi.String("t2.micro"),
			SubnetId:                 subnet1.ID(),
			VpcSecurityGroupIds:      pulumi.StringArray{sg.ID()},
			AssociatePublicIpAddress: pulumi.Bool(true),
			UserData:                pulumi.String(`#!/bin/bash
echo "<html><body><h1>Servidor 1</h1></body></html>" > /var/www/html/index.html
systemctl start nginx`),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("server-1"),
			},
		})
		if err != nil {
			return err
		}

		server2, err := ec2.NewInstance(ctx, "server-2", &ec2.InstanceArgs{
			Ami:                      pulumi.String("ami-0c55b159cbfafe1f0"),
			InstanceType:             pulumi.String("t2.micro"),
			SubnetId:                 subnet2.ID(),
			VpcSecurityGroupIds:      pulumi.StringArray{sg.ID()},
			AssociatePublicIpAddress: pulumi.Bool(true),
			UserData:                pulumi.String(`#!/bin/bash
echo "<html><body><h1>Servidor 2</h1></body></html>" > /var/www/html/index.html
systemctl start nginx`),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("server-2"),
			},
		})
		if err != nil {
			return err
		}

		// Crear balanceador de carga
		lb, err := lb.NewLoadBalancer(ctx, "web-lb", &lb.LoadBalancerArgs{
			SubnetIds:        pulumi.StringArray{subnet1.ID(), subnet2.ID()},
			SecurityGroups:   pulumi.StringArray{sg.ID()},
			LoadBalancerType: pulumi.String("application"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("web-lb"),
			},
		})
		if err != nil {
			return err
		}

		// Exportar información
		ctx.Export("vpcId", vpc.ID())
		ctx.Export("server1PublicIp", server1.PublicIp)
		ctx.Export("server2PublicIp", server2.PublicIp)
		ctx.Export("loadBalancerDns", lb.DnsName)

		return nil
	})
}
```

### Terraform CDK con Go

Terraform CDK permite definir infraestructura utilizando Go:

```go
package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/hashicorp/terraform-cdk-go/cdktf/providers/aws"
)

func NewMyStack(scope constructs.Construct, id string) cdktf.TerraformStack {
	stack := cdktf.NewTerraformStack(scope, &id)

	// Definir proveedor AWS
	aws.NewAwsProvider(stack, jsii.String("AWS"), &aws.AwsProviderConfig{
		Region: jsii.String("us-west-2"),
	})

	// Crear VPC
	vpc := aws.NewVpc(stack, jsii.String("MyVpc"), &aws.VpcConfig{
		CidrBlock: jsii.String("10.0.0.0/16"),
		Tags: &map[string]*string{
			"Name": jsii.String("my-vpc"),
		},
	})

	// Crear subred
	subnet := aws.NewSubnet(stack, jsii.String("MySubnet"), &aws.SubnetConfig{
		VpcId:            vpc.Id(),
		CidrBlock:        jsii.String("10.0.1.0/24"),
		AvailabilityZone: jsii.String("us-west-2a"),
		Tags: &map[string]*string{
			"Name": jsii.String("my-subnet"),
		},
	})

	// Crear grupo de seguridad
	sg := aws.NewSecurityGroup(stack, jsii.String("WebSG"), &aws.SecurityGroupConfig{
		VpcId: vpc.Id(),
		Ingress: &[]aws.SecurityGroupIngress{
			{
				FromPort:   jsii.Number(80),
				ToPort:     jsii.Number(80),
				Protocol:   jsii.String("tcp"),
				CidrBlocks: &[]*string{jsii.String("0.0.0.0/0")},
			},
		},
		Egress: &[]aws.SecurityGroupEgress{
			{
				FromPort:   jsii.Number(0),
				ToPort:     jsii.Number(0),
				Protocol:   jsii.String("-1"),
				CidrBlocks: &[]*string{jsii.String("0.0.0.0/0")},
			},
		},
		Tags: &map[string]*string{
			"Name": jsii.String("web-sg"),
		},
	})

	// Crear instancia EC2
	instance := aws.NewInstance(stack, jsii.String("WebServer"), &aws.InstanceConfig{
		Ami:                      jsii.String("ami-0c55b159cbfafe1f0"),
		InstanceType:             jsii.String("t2.micro"),
		SubnetId:                 subnet.Id(),
		VpcSecurityGroupIds:      &[]*string{sg.Id()},
		AssociatePublicIpAddress: jsii.Bool(true),
		UserData:                jsii.String(`#!/bin/bash
echo "<html><body><h1>Hola desde Terraform CDK</h1></body></html>" > /var/www/html/index.html
systemctl start nginx`),
		Tags: &map[string]*string{
			"Name": jsii.String("web-server"),
		},
	})

	// Exportar información
	cdktf.NewTerraformOutput(stack, jsii.String("public_ip"), &cdktf.TerraformOutputConfig{
		Value: instance.PublicIp(),
	})

	return stack
}

func main() {
	app := cdktf.NewApp(nil)
	NewMyStack(app, "terraform-cdk-go")
	app.Synth()
}
```

## Observabilidad en la Nube

### Logging Estructurado

```go
package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Configurar encoder para JSON
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Crear logger
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encorerConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	logger := zap.New(core).With(
		zap.String("service", "mi-servicio"),
		zap.String("environment", "production"),
	)
	defer logger.Sync()

	// Registrar eventos
	logger.Info("Servicio iniciado",
		zap.String("version", "1.0.0"),
	)

	// Simular procesamiento de solicitud
	logger.Info("Procesando solicitud",
		zap.String("request_id", "abc-123"),
		zap.String("user_id", "user-456"),
	)

	// Registrar métricas
	tiempoInicio := time.Now()
	time.Sleep(100 * time.Millisecond) // Simular procesamiento
	logger.Info("Solicitud completada",
		zap.String("request_id", "abc-123"),
		zap.Duration("duration_ms", time.Since(tiempoInicio)),
		zap.Int("status_code", 200),
	)

	// Registrar error
	logger.Error("Error al conectar con base de datos",
		zap.String("db_host", "db.example.com"),
		zap.Int("retry_count", 3),
		zap.Error(errors.New("timeout exceeded")),
	)
}
```

### Métricas con Prometheus

```go
package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Contador de solicitudes totales
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "Número total de solicitudes procesadas",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Histograma de tiempos de respuesta
	responseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_response_time_seconds",
			Help:    "Tiempo de respuesta en segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Gauge para conexiones activas
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_active_connections",
			Help: "Número actual de conexiones activas",
		},
	)
)

func main() {
	// Endpoint para métricas de Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Endpoint de ejemplo
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		// Incrementar contador de conexiones activas
		activeConnections.Inc()
		defer activeConnections.Dec()

		// Medir tiempo de respuesta
		timer := prometheus.NewTimer(responseTime.With(prometheus.Labels{
			"method":   r.Method,
			"endpoint": "/api/data",
		}))
		defer timer.ObserveDuration()

		// Simular procesamiento
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		// Simular diferentes códigos de estado
		statusCode := 200
		if rand.Float64() < 0.1 {
			statusCode = 500
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error interno del servidor"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Datos procesados correctamente"))
		}

		// Incrementar contador de solicitudes
		requestCounter.With(prometheus.Labels{
			"method":   r.Method,
			"endpoint": "/api/data",
			"status":   fmt.Sprintf("%d", statusCode),
		}).Inc()
	})

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Println("Métricas disponibles en http://localhost:8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Trazabilidad con OpenTelemetry

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

// Inicializar trazabilidad
func initTracer() func() {
	ctx := context.Background()

	// Crear conexión con el colector OTLP
	conn, err := grpc.DialContext(ctx, "localhost:4317", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar con colector OTLP: %v", err)
	}

	// Crear exportador OTLP
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithGRPCConn(conn),
		),
	)
	if err != nil {
		log.Fatalf("Error al crear exportador OTLP: %v", err)
	}

	// Crear recurso con metadatos del servicio
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("mi-servicio"),
			semconv.ServiceVersionKey.String("1.0.0"),
			attribute.String("environment", "production"),
		),
	)
	if err != nil {
		log.Fatalf("Error al crear recurso: %v", err)
	}

	// Crear proveedor de trazas
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	// Configurar propagación de contexto
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Devolver función de limpieza
	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Error al cerrar proveedor de trazas: %v", err)
		}
		conn.Close()
	}
}

// Servicio de usuarios (simulado)
func userService(ctx context.Context, userID string) (string, error) {
	tracer := otel.Tracer("user-service")
	ctx, span := tracer.Start(ctx, "GetUserDetails")
	defer span.End()

	span.SetAttributes(attribute.String("user.id", userID))

	// Simular procesamiento
	time.Sleep(50 * time.Millisecond)

	return fmt.Sprintf("Usuario %s", userID), nil
}

// Servicio de productos (simulado)
func productService(ctx context.Context, productID string) (string, error) {
	tracer := otel.Tracer("product-service")
	ctx, span := tracer.Start(ctx, "GetProductDetails")
	defer span.End()

	span.SetAttributes(attribute.String("product.id", productID))

	// Simular procesamiento
	time.Sleep(30 * time.Millisecond)

	return fmt.Sprintf("Producto %s", productID), nil
}

// Manejador HTTP para pedidos
func orderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("order-service")

	// Extraer parámetros
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		productID = "unknown"
	}

	// Crear span para procesar pedido
	ctx, span := tracer.Start(ctx, "ProcessOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("order.user_id", userID),
		attribute.String("order.product_id", productID),
	)

	// Obtener detalles del usuario
	userDetails, err := userService(ctx, userID)
	if err != nil {
		http.Error(w, "Error al obtener detalles del usuario", http.StatusInternalServerError)
		return
	}

	// Obtener detalles del producto
	productDetails, err := productService(ctx, productID)
	if err != nil {
		http.Error(w, "Error al obtener detalles del producto", http.StatusInternalServerError)
		return
	}

	// Simular procesamiento del pedido
	time.Sleep(20 * time.Millisecond)

	// Responder
	response := fmt.Sprintf("Pedido procesado para %s, producto: %s", userDetails, productDetails)
	io.WriteString(w, response)
}

func main() {
	// Inicializar trazabilidad
	cleanup := initTracer()
	defer cleanup()

	// Configurar rutas con instrumentación
	http.Handle("/api/orders", otelhttp.NewHandler(
		http.HandlerFunc(orderHandler),
		"orders-endpoint",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	))

	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Ejercicios Prácticos

### Ejercicio 1: Microservicio en Contenedor

Desarrolla un microservicio en Go que proporcione una API REST para gestionar una lista de tareas, y despliégalo en un contenedor Docker.

1. Implementa endpoints CRUD para tareas.
2. Utiliza una base de datos (puede ser en memoria para simplificar).
3. Añade logging estructurado.
4. Crea un Dockerfile multi-etapa.
5. Configura la aplicación mediante variables de entorno.

### Ejercicio 2: Función Serverless

Crea una función serverless en Go que procese imágenes subidas a un bucket de almacenamiento en la nube.

1. Implementa la función para AWS Lambda o Google Cloud Functions.
2. Configura los triggers para activar la función cuando se suban imágenes.
3. Procesa las imágenes (redimensionar, convertir formato, etc.).
4. Guarda las imágenes procesadas en otro bucket.
5. Implementa manejo de errores y reintentos.

### Ejercicio 3: Aplicación Kubernetes

Desarrolla una aplicación compuesta por múltiples microservicios y despliégala en Kubernetes.

1. Crea al menos dos microservicios que se comuniquen entre sí.
2. Implementa patrones de resiliencia (circuit breaker, retry).
3. Configura observabilidad (métricas, logs, trazas).
4. Define manifiestos de Kubernetes para el despliegue.
5. Implementa pruebas de integración.

## Conclusiones

Go se ha establecido como un lenguaje de primera clase para el desarrollo de aplicaciones en la nube debido a su rendimiento, simplicidad y excelente soporte para la concurrencia. Sus características lo hacen ideal para microservicios, aplicaciones serverless, herramientas de infraestructura y operadores de Kubernetes.

Al desarrollar aplicaciones cloud-native con Go, es importante seguir los principios de diseño para la nube, implementar patrones de resiliencia y configurar adecuadamente la observabilidad para garantizar que las aplicaciones sean robustas, escalables y fáciles de mantener.

La combinación de Go con tecnologías como Docker, Kubernetes, y las plataformas serverless de los principales proveedores de nube, proporciona un ecosistema potente para construir sistemas distribuidos modernos.

## Referencias

1. Go en la Nube: https://go.dev/solutions/cloud/
2. Docker con Go: https://docs.docker.com/language/golang/
3. Kubernetes Client-Go: https://github.com/kubernetes/client-go
4. AWS SDK para Go: https://aws.github.io/aws-sdk-go-v2/docs/
5. Google Cloud SDK para Go: https://cloud.google.com/go/docs
6. Azure SDK para Go: https://github.com/Azure/azure-sdk-for-go
7. Pulumi con Go: https://www.pulumi.com/docs/intro/languages/go/
8. Terraform CDK con Go: https://developer.hashicorp.com/terraform/cdktf
9. OpenTelemetry para Go: https://opentelemetry.io/docs/go/
10. "Cloud Native Go" - Matthew A. Titmus
11. "Go Programming Blueprints" - Mat Ryer
12. "Building Microservices with Go" - Nic Jackson