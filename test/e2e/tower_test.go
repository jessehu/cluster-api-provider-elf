package e2e

import (
	"flag"
	"os"
	"strconv"

	. "github.com/onsi/gomega"

	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	infrav1 "github.com/smartxworks/cluster-api-provider-elf/api/v1alpha4"
	"github.com/smartxworks/cluster-api-provider-elf/pkg/service"
)

var (
	elfTemplate   = os.Getenv("ELF_TEMPLATE")
	towerUsername = os.Getenv("TOWER_USERNAME")
	towerPassword = os.Getenv("TOWER_PASSWORD")

	towerServer     string
	towerServerPort int
	vmService       service.VMService
)

func init() {
	flag.StringVar(&towerServer, "e2e.towerServer", os.Getenv("TOWER_SERVER"), "the tower server used for e2e tests")
	port, _ := strconv.Atoi(os.Getenv("TOWER_SERVER_PORT"))
	flag.IntVar(&towerServerPort, "e2e.towerServerPort", port, "the tower server port used for e2e tests")
}

func initElfSession() {
	var err error
	vmService, err = service.NewVMService(infrav1.Tower{
		Server:   towerServer,
		Port:     towerServerPort,
		Username: towerUsername,
		Password: towerPassword}, ctrllog.Log)
	Expect(err).ShouldNot(HaveOccurred())

	template, err := vmService.GetVMTemplate(elfTemplate)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(template.LocalID).ShouldNot(BeEmpty())
}
