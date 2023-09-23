package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	pb "addressbook/addressbookpb"

	"google.golang.org/protobuf/proto"
)

func promptForAddress(r io.Reader, book *pb.AddressBook) error {
	var (
		id             int32
		selectedPerson *pb.Person
	)

	reader := bufio.NewReader(r)

	fmt.Print("Enter person ID number: ")

	n, err := fmt.Fscanf(reader, "%d\n", &id)
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("You can only update one address at a time.")
	}

	for _, address := range book.People {
		if address.Id == id {
			selectedPerson = address
			break
		}
	}

	if selectedPerson == nil {
		return fmt.Errorf("No person is found.")
	}

	fmt.Println(selectedPerson.String())

	fmt.Printf("Enter new person's name (%s), press enter to skip: ", selectedPerson.Name)
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)
	if name != "" {
		selectedPerson.Name = name
	}

	fmt.Printf("Enter new person's email (%s), press enter to skip: ", selectedPerson.Email)
	email, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if email != "" {
		selectedPerson.Email = strings.TrimSpace(email)
	}

	for _, phone := range selectedPerson.Phones {
		fmt.Printf("Enter a new phone number (%s), press enter to skip: ", phone.Number)
		pNumber, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		pNumber = strings.TrimSpace(pNumber)
		if pNumber != "" {
			phone.Number = pNumber
		}

		fmt.Printf("Is this a mobile, home, or work phone? (%s) press enter to skip: ", phone.Type)
		phoneType, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		phoneType = strings.TrimSpace(phoneType)

		switch phoneType {
		case "mobile":
			phone.Type = pb.Person_PHONETYPE_MOBILE
		case "home":
			phone.Type = pb.Person_PHONETYPE_HOME
		case "work":
			phone.Type = pb.Person_PHONETYPE_WORK
		case "":

		default:
			fmt.Printf("Unknown phone type %q. Using default.\n", phoneType)
		}

	}

	return nil
}

func getFileName() (string, error) {
	// Only accept filename as the argument
	// program name + filename = 2 arguments
	if len(os.Args) != 2 {
		return "", errors.New("Invalid arguments, only one filename is accepted ")
	}

	return os.Args[1], nil
}

func readFile(fname string) ([]byte, error) {
	in, err := os.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File not found. Creating new file after address is input.\n", fname)
			return []byte{}, nil
		} else {
			return nil, fmt.Errorf("Error reading file: %s", fname)
		}
	}
	return in, nil
}

func run() error {
	fname, err := getFileName()
	if err != nil {
		return err
	}

	in, err := readFile(fname)
	if err != nil {
		return err
	}

	book := &pb.AddressBook{}
	if err := proto.Unmarshal(in, book); err != nil {
		return fmt.Errorf("Failed to parse address book: %s", err)
	}

	err = promptForAddress(os.Stdin, book)
	if err != nil {
		return fmt.Errorf("Error with address: %s", err)
	}

	out, err := proto.Marshal(book)
	if err != nil {
		return fmt.Errorf("Failed to encode address book: %s", err)
	}

	if err := os.WriteFile(fname, out, 0644); err != nil {
		return fmt.Errorf("Failed to write address book: %s", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
