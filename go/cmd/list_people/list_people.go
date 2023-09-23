package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	pb "addressbook/addressbookpb"
	"google.golang.org/protobuf/proto"
)

func writePerson(w io.Writer, p *pb.Person) {
	fmt.Fprintln(w, "Person ID:", p.Id)
	fmt.Fprintln(w, "  Name:", p.Name)
	if p.Email != "" {
		fmt.Fprintln(w, "  E-mail address:", p.Email)
	}

	for _, pn := range p.Phones {
		switch pn.Type {
		case pb.Person_PHONETYPE_MOBILE:
			fmt.Fprintf(w, "  Mobile phone #: ")
		case pb.Person_PHONETYPE_WORK:
			fmt.Fprintf(w, "  Work phone #: ")
		case pb.Person_PHONETYPE_HOME:
			fmt.Fprintf(w, "  Home phone #: ")
		}
		fmt.Fprintln(w, pn.Number)
	}
}

func getFileName() (string, error) {
	// Only accept filename as the argument
	// program name + filename = 2 arguments
	if len(os.Args) != 2 {
		return "", errors.New("Invalid arguments, only one filename is accepted ")
	}

	return os.Args[1], nil
}

func listPeople(w io.Writer, book *pb.AddressBook) {
	for _, p := range book.People {
		writePerson(w, p)
	}
}

func run() error {
	fname, err := getFileName()
	if err != nil {
		return err
	}

	// [START unmarshal_proto]
	// Read the existing address book
	in, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Error reading file: %w", err)
	}

	book := &pb.AddressBook{}
	if err := proto.Unmarshal(in, book); err != nil {
		return fmt.Errorf("Failed to parse address book: %s", err)
	}

	listPeople(os.Stdout, book)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
