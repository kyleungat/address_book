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

func promptForAddress(r io.Reader) (*pb.Person, error) {
	p := &pb.Person{}

	reader := bufio.NewReader(r)

	fmt.Print("Enter person ID number: ")

	if _, err := fmt.Fscanf(reader, "%d\n", &p.Id); err != nil {
		return p, err
	}

	fmt.Print("Enter name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return p, err
	}

	p.Name = strings.TrimSpace(name)

	fmt.Print("Enter email address (blank for none): ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return p, err
	}
	p.Email = strings.TrimSpace(email)

	for {
		fmt.Print("Enter a phone number (or leave blank to finish): ")
		phone, err := reader.ReadString('\n')
		if err != nil {
			return p, err
		}
		phone = strings.TrimSpace(phone)
		if phone == "" {
			break
		}

		pn := &pb.Person_PhoneNumber{
			Number: phone,
		}

		fmt.Print("Is this a mobile, home, or work phone? ")
		phoneType, err := reader.ReadString('\n')
		if err != nil {
			return p, err
		}
		phoneType = strings.TrimSpace(phoneType)

		switch phoneType {
		case "mobile":
			pn.Type = pb.Person_PHONETYPE_MOBILE
		case "home":
			pn.Type = pb.Person_PHONETYPE_HOME
		case "work":
			pn.Type = pb.Person_PHONETYPE_WORK
		default:
			fmt.Printf("Unknown phone type %q. Using default.\n", phoneType)
		}

		p.Phones = append(p.Phones, pn)
	}

	return p, nil
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

	addr, err := promptForAddress(os.Stdin)
	if err != nil {
		return fmt.Errorf("Error with address: %s", err)
	}

	book.People = append(book.People, addr)

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
