
# cloudfirestore
`cloudfirestore` is a further simplification of the Cloud Firestore Go SDK library.

- [Firestore(GoogleCloud)](https://cloud.google.com/firestore?hl=ja)

## firestore SDK

```
type Book struct {
  Title  string
  Author string
}

client, _ := firestore.NewClient(ctx, firestore.DetectProjectID)

b := &Book{}
docID := "xxx"
client.Collction("books").Doc(docID).Get(ctx, docID)
```

## cloudfirestore
```
type Book struct {
  ID     string
  Title  string
  Author string
}

func (b Book) GenerateKey(_ context.Context) string {
  return "books/" + b.ID
}

c, _ := cloudfirestore.New(ctx)

b := &Book{ID:"xxx"}
c.Get(ctx, b)
```



## Note
- The data must be capable of generating a document path by itself.