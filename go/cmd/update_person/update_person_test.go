package main

import (
	pb "addressbook/addressbookpb"
	"google.golang.org/protobuf/proto"
	"strings"
	"testing"
)

func TestPromptForAddress(t *testing.T) {
	in := `1
	
	123123@gmail.com
	
	home
	`

	book := &pb.AddressBook{
		People: []*pb.Person{
			{
				Id:    1,
				Name:  "123",
				Email: "123@gmail.com",
				Phones: []*pb.Person_PhoneNumber{
					{
						Number: "123456789",
						Type:   pb.Person_PHONETYPE_MOBILE,
					},
				},
			},
			{
				Id:    2,
				Name:  "234",
				Email: "234@gmail.com",
			},
			{
				Id:    3,
				Name:  "345",
				Email: "345@gmail.com",
			},
		},
	}

	err := promptForAddress(strings.NewReader(in), book)
	if err != nil {
		t.Fatalf("promptForDelete(%q) had unexpected error: %s", in, err.Error())
	}

	want := []*pb.Person{
		{
			Id:    1,
			Name:  "123",
			Email: "123123@gmail.com",
			Phones: []*pb.Person_PhoneNumber{
				{
					Number: "123456789",
					Type:   pb.Person_PHONETYPE_HOME,
				},
			},
		},
		{
			Id:    2,
			Name:  "234",
			Email: "234@gmail.com",
		},
		{
			Id:    3,
			Name:  "345",
			Email: "345@gmail.com",
		},
	}

	if len(book.People) != len(want) {
		t.Errorf("want %d addresses, got %d", len(want), len(book.People))
	}

	for i := 0; i < len(book.People); i++ {
		if !proto.Equal(book.People[i], want[i]) {
			t.Errorf("want address %q, got %q", want[i], book.People[i])
		}

	}
}
