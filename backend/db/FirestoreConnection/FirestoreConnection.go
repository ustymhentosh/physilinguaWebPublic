package firestoreconnection

import (
	"context"
	"fmt"
	"slices"
	"strings"
	prbl "veles/db/Problem"

	"cloud.google.com/go/firestore"
)

var (
	ctx = context.Background()
)

type FirestoreConnection struct {
	client *firestore.Client
}

func New(client *firestore.Client) FirestoreConnection {
	fc := FirestoreConnection{client: client}
	return fc
}

type SimpleOutput struct {
	Ready   []string `firestore:"ready"`
	Pending []string `firestore:"pending"`
}

type OneProblemOutput struct {
	Ready []string `firestore:"ready"`
}

type Statuses struct {
	Subs map[string][]string `firestore:"-"`
}

// Gets list of all problems present in database, returns slice  []Probmem
func (fc FirestoreConnection) GetProblemsList() []prbl.Problem {
	// getting to collection
	problems_ref := fc.client.Collection("Overview")
	updetes_doc := problems_ref.Doc("simple")
	updates_ref, _ := updetes_doc.Get(ctx)

	var updates SimpleOutput
	err := updates_ref.DataTo(&updates)

	if err != nil {
		fmt.Println("unsuccesfull parsing Simple doc", err)
	}

	var all_docs_ids []prbl.Problem

	// Unpacking all documents into []Problem structure for further use
	for _, val := range updates.Pending {
		all_docs_ids = append(all_docs_ids, prbl.Problem{
			Id:     strings.Replace(fmt.Sprintf("%v", val), "-", ".", -1),
			Status: "inProgress",
		})
	}
	for _, val := range updates.Ready {
		all_docs_ids = append(all_docs_ids, prbl.Problem{
			Id:     strings.Replace(fmt.Sprintf("%v", val), "-", ".", -1),
			Status: "done",
		})
	}

	return all_docs_ids
}

// Saves given data to document of a problem, if problem document
// did not exists, -> creates new document
func (fc FirestoreConnection) SaveProblemSubmission(number string,
	text string, answer string, extensions []string, comment string) []string {
	// Get collection and get all documents,
	problems_ref := fc.client.Collection("FromUsers")
	all_docs, _ := problems_ref.DocumentRefs(ctx).GetAll()
	added := false

	var images_names []string

	// If document exists -> update
	for _, i := range all_docs {
		id := i.ID
		if strings.ReplaceAll(id, "-", ".") == number {
			content_ref, _ := i.Get(ctx)
			subs := content_ref.Data()
			new_sub := len(subs) + 1
			images_names = getNumsForImages(id, extensions, new_sub)
			updateEntries(i, number, text, answer, images_names, comment, new_sub)
			fc.logUpdateToAdmin(number, fmt.Sprint(new_sub))
			added = true
		}
	}

	// Create new document if not exists
	if !added {
		images_names = getNumsForImages(strings.ReplaceAll(number, ".", "-"), extensions, 1)
		createDocument(problems_ref, number, text, answer, images_names, comment)
		fc.logUpdateToAdmin(number, "1")
	}

	return images_names
}

func updateEntries(cur_doc *firestore.DocumentRef, number string,
	text string, answer string, images []string, comment string, new_sub int) {
	cur_doc.Update(ctx, []firestore.Update{{Path: fmt.Sprint(new_sub) + ".text", Value: text}})
	cur_doc.Update(ctx, []firestore.Update{{Path: fmt.Sprint(new_sub) + ".answer", Value: answer}})
	cur_doc.Update(ctx, []firestore.Update{{Path: fmt.Sprint(new_sub) + ".images", Value: images}})
	cur_doc.Update(ctx, []firestore.Update{{Path: fmt.Sprint(new_sub) + ".comment", Value: comment}})
}

func createDocument(collection_ref *firestore.CollectionRef, number string,
	text string, answer string, images []string, comment string) {
	collection_ref.Doc(strings.ReplaceAll(number, ".", "-")).Set(ctx, map[string]interface{}{
		"1": map[string]interface{}{
			"text":    text,
			"answer":  answer,
			"images":  images,
			"comment": comment,
		},
	})
}

func (fc FirestoreConnection) logUpdateToAdmin(number_of_problem string, num_of_sub string) {
	dashed_num := strings.ReplaceAll(number_of_problem, ".", "-")
	admin_ref := fc.client.Collection("Overview")
	doc_ref, _ := admin_ref.Doc("submissions-status").Get(ctx)
	simple_ref, _ := admin_ref.Doc("simple").Get(ctx)
	problems := fc.GetProblemsList()
	ready_problems := prbl.GetFilteredIdList(problems, "done")
	pending_problems := prbl.GetFilteredIdList(problems, "inProgress")

	if !slices.Contains(ready_problems, number_of_problem) {
		pending_problems = append(pending_problems, number_of_problem)
		simple_ref.Ref.Update(ctx, []firestore.Update{{
			Path:  "pending",
			Value: pending_problems}})
	}

	doc_ref.Ref.Update(ctx, []firestore.Update{{
		Path:  dashed_num + "." + fmt.Sprint(num_of_sub),
		Value: "not-viewed"}})
}

func getNumsForImages(number string, q []string, sub int) []string {
	var names []string
	for i := 1; i <= len(q); i++ {
		names = append(names, number+fmt.Sprintf("(%d)", sub)+fmt.Sprintf("(%d)", i))
	}
	return names
}

// Checks whether given problem is full
func (fc FirestoreConnection) GetNumsOfSubmissions(number string) int {
	overview_ref := fc.client.Collection("Overview")
	subms_doc := overview_ref.Doc("submissions-status")
	subms_ref, _ := subms_doc.Get(ctx)
	dataMap := subms_ref.Data()
	innerMap, _ := dataMap[strings.ReplaceAll(number, ".", "-")].(map[string]interface{})
	length := len(innerMap)
	return length
}
