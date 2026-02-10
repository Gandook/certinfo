package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Gandook/certinfo/certinfo"
)

// printUsageGuide prints a usage guide to stderr.
func printUsageGuide() {
	_, err := fmt.Fprintf(os.Stderr, `Usage: certinfo-cli <command> [options]

Commands:
	certinfo	Show the useful information about a certain DigSig X.509 certificate
`)
	if err != nil {
		return
	}
}

// runCertInfo executes a "certinfo" command to retrieve the useful information about
// a certain DigSig X.509 certificate.
func runCertInfo(retriever certinfo.InfoRetriever, args []string) error {
	command := flag.NewFlagSet("certinfo", flag.ExitOnError)
	daid := command.String("daid", "", "The certificate's DAID")
	cid := command.String("cid", "", "The certificate's CID")
	err := command.Parse(args)
	if err != nil {
		return err
	}

	info, retrieveErr := retriever.Retrieve(*daid, *cid)
	if retrieveErr != nil {
		return retrieveErr
	}

	fmt.Printf("CID: %s\nDAID: %s\nIssuer: %s\nSubject: %s\nNotBefore: %s\nNotAfter: %s\n",
		info.CID,
		info.DAID,
		info.Issuer,
		info.Subject,
		info.NotBefore,
		info.NotAfter)

	return nil
}

func main() {
	retriever := certinfo.NewRetriever()

	// Printing a usage guide if the command has too few arguments.
	if len(os.Args) < 5 {
		printUsageGuide()
		os.Exit(1)
	}

	if os.Args[1] == "certinfo" {
		if err := runCertInfo(retriever, os.Args[2:]); err != nil {
			log.Fatalf("Error in certinfo command: %v", err)
		}
	} else {
		printUsageGuide()
		os.Exit(1)
	}
}
