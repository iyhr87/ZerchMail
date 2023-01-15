package main

import "indexer/mymailstruct"

func main() {
	root := "/home/ivan/zinc/data/enron_mail_20110402"
	mymailstruct.IndexInZincSearc(root)
}
