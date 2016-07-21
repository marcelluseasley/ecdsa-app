package main

import (
	"crypto/ecdsa" //main cryptographic functions
	"crypto/elliptic" //used to get the elliptic curve
	"crypto/rand"     //Package rand implements a cryptographically secure pseudorandom number generator.
	"crypto/sha256" //used for the sha256 hash and digest
	"crypto/x509" //used for key parsing and mashalling
	"encoding/base64" //used to convert signature to base64
	"encoding/pem" //used to create the PEM format string
	"fmt" //console output
	"math/big" //used to store very large numbers (bytes in this case)
	"os" //operating system interaction
)

/*


Using a language of your choice, provide an application that meets the following requirements:

Given a string input of up to 250 characters, return a JSON response compliant to the schema defined below.

You are responsible for generating a public/private ECDSA keypair and persisting the keypair on the filesystem
Subsequent invocations of your application should read from the same files

Document your code, at a minimum defining parameter types and return values for any public methods

Include Unit Test(s) with instructions on how a Continuous Integration system can execute your test(s)




*/

func main() {

	//if not two arguments print correct usage and exit
	if len(os.Args) != 2{
		fmt.Printf("Usage: %s (input string)\n", os.Args[0])
		os.Exit(1)
	}

	//PrivateKey represents a ECDSA private key.
	privateKey := new(ecdsa.PrivateKey)

	//needed to generate the private key; P256 returns a Curve which implements P-256 (see FIPS 186-3, section D.2.3)
	ellipticCurve := elliptic.P256()

	//using the elliptic curve and random number, generate the private key
	privateKey, err := ecdsa.GenerateKey(ellipticCurve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//the privatekey struct contains an instance of the public key associated with itself (the private key)
	//
	//
	// publicKey := privateKey.PublicKey
	/*
	fmt.Println("Private key: ")
	fmt.Printf("%x \n", privateKey)

	fmt.Println("Public key: ")
	fmt.Printf("%x \n", publicKey)
	*/
	//TODO: write these keys to filesystem
	/*

		After private and public key, need to make Base64 encoded cryptographic signature
		using private key and SHA256 digest of the input
	*/

	//myHash is the digest returned by the sha256 hash function
	//Sum256 returns the SHA256 checksum of the data.
	myHash := sha256.Sum256([]byte(os.Args[1]))

	//needed by the ecdsa.Sign function which returns the sig in two different variables
	signaturePart1 := big.NewInt(0)
	signaturePart2 := big.NewInt(0)

	//returns signature parts; function takes random number, private key, and the hash digest
	signaturePart1, signaturePart2, err = ecdsa.Sign(rand.Reader, privateKey, myHash[:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//construct the full signature
	signature := signaturePart1.Bytes()
	signature = append(signature, signaturePart2.Bytes()...)

	fmt.Printf("Signature: %s\n", base64.StdEncoding.EncodeToString(signature))

	//verify
	//verifystatus := ecdsa.Verify(&publicKey, myHash[:], signaturePart1, signaturePart2)
	//fmt.Println(verifystatus)

	//MarshalECPrivateKey marshals an EC private key into ASN.1, DER format.
	asn_der, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var pp = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn_der, // The decoded bytes of the contents. Typically a DER encoded ASN.1 structure.
	}

	fmt.Printf("%s \n", pem.EncodeToMemory(pp))
}
