package techpuzzle7

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	documentai "cloud.google.com/go/documentai/apiv1"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	documentaipb "google.golang.org/genproto/googleapis/cloud/documentai/v1"
)

var (
	storageClient   *storage.Client
	firestoreClient *firestore.Client
)

const processor string = "projects/872946784243/locations/us/processors/9191950ebef6ca96"
const projectID string = "techpuzzle7"

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

type Invoice struct {
	Invoice_id string
	Amount     string
	Paid       bool
}

// ProcessInvoice is executed when a file is uploaded to the Cloud Storage bucket you
// created for uploading invoices. It processes the invoice for text and updates the
// invoice data in Firebase.
func ProcessInvoice(ctx context.Context, event GCSEvent) error {

	// Read the file from GCS that triggered the event and get it as a byte array
	if event.Bucket == "" {
		return fmt.Errorf("empty file.Bucket")
	}
	if event.Name == "" {
		return fmt.Errorf("empty file.Name")
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	bucket := storageClient.Bucket(event.Bucket)

	rc, err := bucket.Object(event.Name).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("readFile: unable to open file from bucket %q, file %q: %v", event.Bucket, event.Name, err)
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("readFile: unable to read data from bucket %q, file %q: %v", event.Bucket, event.Name, err)
	}

	//Detect the mime type of the file
	mimeType := http.DetectContentType(data)

	//Invoke the document processor
	client, error := documentai.NewDocumentProcessorClient(ctx)
	if error != nil {
		return fmt.Errorf("failed to create new document processor client: %v", error)
	}
	defer client.Close()
	req := &documentaipb.ProcessRequest{
		Source: &documentaipb.ProcessRequest_RawDocument{
			RawDocument: &documentaipb.RawDocument{
				Content:  data,
				MimeType: mimeType,
			},
		},
		Name: processor,
	}
	resp, err := client.ProcessDocument(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to process response: %v", err)
	}

	// Parse through the detected fields on the invoice
	document := resp.GetDocument()
	entities := document.GetEntities()
	var inv Invoice
	for _, entity := range entities {
		if entity.GetType() == "invoice_id" {
			inv.Invoice_id = entity.GetMentionText()
		}
		if entity.GetType() == "total_amount" {
			sAmount := entity.GetMentionText()
			inv.Amount = sAmount
		}
	}
	fmt.Printf("Invoice ID: %s\n", inv.Invoice_id)
	fmt.Printf("Total Amount: %s\n", inv.Amount)

	// Get a Firestore client.
	firestoreClient, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to create firestore.NewClient: %v", err)
	}
	defer firestoreClient.Close()

	//Look for the invoice in Firestore
	iter := firestoreClient.Collection("invoices").Where("invoice_id", "==", inv.Invoice_id).Documents(ctx)

	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Invoice %s does not exist.\n", inv.Invoice_id)
			// TODO: Implement logic if invoice doesn't exist
			return fmt.Errorf("failed to convert amount: %v", err)
		}

		// Taking care of cases where the invoice is already paid or the amount doesn't match the expected amount
		invoiceMap := dsnap.Data()
		var invData Invoice
		invData.Invoice_id = invoiceMap["invoice_id"].(string)
		invData.Paid = invoiceMap["paid"].(bool)
		if invData.Paid {
			fmt.Printf("The invoice %s has already been paid!\n", invData.Invoice_id)
			// TODO: Implement duplicate invoice logic
		}

		// Update the Firestore record
		_, err = dsnap.Ref.Update(ctx, []firestore.Update{
			{
				Path:  "paid",
				Value: true,
			}})

		if err != nil {
			return fmt.Errorf("failed to update invoice: %v", err)
		}

	}

	return nil
}
