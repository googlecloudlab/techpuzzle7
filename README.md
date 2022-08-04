# techpuzzle7

**The story ...**
Your company processes payments received by mail.  When a payment is received, the invoice is scanned and uploaded as a file into a Google Cloud Storage bucket.  Today, the invoices are examined by staff and the corresponding order, which is contained in a document database, is updated to indicate that the payment was received.  You want to reduce the toil on your staff and reduce the latency from when payment is received to when the order is flagged as paid.

**The puzzle ...**
Your Google Field Sales Rep did a great demonstration of a product called Document AI and you think this may be just what you need.  You want to employ Document AI to handle new documents arriving in a GCS bucket and update a corresponding Firestore database.  You think you will best achieve this by triggering a Cloud Function when new invoice images arrive in the bucket, processing the image with Document AI and then updating the database with the results.

**Notes ...**
Can you handle the case where a duplicate invoice is received?  What if the payment amount of the invoice doesn't match the expected amount?  What if we receive an invoice and don't know which order to apply it to?   None of these are required to submit a solution.  Merely document processing and the corresponding database update will be sufficient.

**Artifacts ...**
A set of (five) sample invoice images are supplied as well as an export from a Firestore database that contains the orders corresponding to each of the invoices.


**Solution***
[Techpuzzle 7 Solution](Techpuzzle7-Solution.pdf)
