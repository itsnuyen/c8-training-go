package internal

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsnuyen/go-zebee-worker/internal/camunda"
)

var camundaClient *camunda.Client

func init() {
	log.Println("Initializing app Config")
	InitViper()

	go InitCamundaZebeeProcess()
}

func InitCamundaZebeeProcess() {

	shutdownZebee := make(chan bool, 1)
	SetupShutdownBarrier(shutdownZebee)
	cClient := camunda.Client{}
	client, err := cClient.InitZeebeClient(CamundaClientId, CamundaClientSecret, CamundaGatewayURL)
	if err != nil {
		panic(err)
	}
	camundaClient = client

	// Deploy Process
	// comment out if not wanted
	//camundaClient.DeployFormDefinition(camunda.OrderProcessBpmn)
	//camundaClient.DeployFormDefinition(camunda.CorrectDataForm)
	//camundaClient.DeployFormDefinition(camunda.OrderDiscountDmn)
	//camundaClient.DeployFormDefinition(camunda.PaymentProcessBpmn)

	defer camundaClient.CloseZeebeClient()
	deductCreditWorker := camundaClient.StartFetchDeductCreditWorker()
	chargeCreditWorker := camundaClient.StartFetchChargeCreditWorker()
	paymentWorker := camundaClient.StartFetchPaymentDonetWorker()
	invokePaymentWorker := camundaClient.StartFetchInvokePaymentWorker()
	<-shutdownZebee
	chargeCreditWorker.AwaitClose()
	deductCreditWorker.AwaitClose()
	paymentWorker.AwaitClose()
	invokePaymentWorker.AwaitClose()
}

func SetupShutdownBarrier(done chan bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
}
