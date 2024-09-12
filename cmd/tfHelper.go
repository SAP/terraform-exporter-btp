package cmd

import (
	"context"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/hashicorp/terraform-exec/tfexec"
)

// Constants for TF version for Terraform providers e.g. for SAP BTP
const BtpProviderVersion = "v1.6.0"

func GetTfStateData(configDir string) ([]byte, error) {
	execPath, err := exec.LookPath("terraform")
	if err != nil {
		log.Fatalf("error finding Terraform: %v", err)
		return nil, err
	}
	// create a new Terraform instance
	tf, err := tfexec.NewTerraform(configDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %v", err)
		return nil, err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %v", err)
		return nil, err
	}
	err = tf.Apply(context.Background())
	if err != nil {
		log.Fatalf("error running Apply: %v", err)
		return nil, err
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %v", err)
		return nil, err
	}

	jsonBytes, err := json.Marshal(state.Values.RootModule.Resources[0].AttributeValues)
	if err != nil {
		log.Fatalf("error json.Marshal: %v", err)
		return nil, err
	}

	return jsonBytes, nil

}
