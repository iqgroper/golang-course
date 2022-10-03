package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type ServerTestCase struct {
	TestDescribtion string
	URL             string
	AccessToken     string
	Query           string
	OrderField      string
	OrderBy         int
	Response        string
	StatusCode      int
}

func TestSearchServer(t *testing.T) {
	cases := []ServerTestCase{
		{
			TestDescribtion: "OrderBy checks",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "Age",
			OrderBy:         1,
			Response:        `[{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "OrderBy checks",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "Age",
			OrderBy:         0,
			Response:        `[{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "OrderBy checks",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "Age",
			OrderBy:         -1,
			Response:        `[{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "OrderField checks",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "Id",
			OrderBy:         1,
			Response:        `[{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "OrderField checks",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "Id",
			OrderBy:         -1, // для наглядности, при нуле будет так же как при единице при OrderField=Id
			Response:        `[{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"},{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "OrderField checks",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "Age",
			OrderBy:         1,
			Response:        `[{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "OrderField checks",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "Name",
			OrderBy:         1,
			Response:        `[{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "OrderField checks",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "Name",
			OrderBy:         -1,
			Response:        `[{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "Blank OrderField check",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"RoseCarney\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField:      "",
			OrderBy:         1,
			Response:        `[{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"}]`,
			StatusCode:      http.StatusOK,
		},
		{
			TestDescribtion: "Blank query check",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "",
			OrderField:      "Age",
			OrderBy:         1,
			Response:        `[{"ID":1,"Name":"HildaMayer","Age":21,"About":"Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n","Gender":"female"},{"ID":15,"Name":"AllisonValdez","Age":21,"About":"Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n","Gender":"male"},{"ID":23,"Name":"GatesSpencer","Age":21,"About":"Dolore magna magna commodo irure. Proident culpa nisi veniam excepteur sunt qui et laborum tempor. Qui proident Lorem commodo dolore ipsum.\n","Gender":"male"},{"ID":0,"Name":"BoydWolf","Age":22,"About":"Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n","Gender":"male"},{"ID":14,"Name":"NicholsonNewman","Age":23,"About":"Tempor minim reprehenderit dolore et ad. Irure id fugiat incididunt do amet veniam ex consequat. Quis ad ipsum excepteur eiusmod mollit nulla amet velit quis duis ut irure.\n","Gender":"male"},{"ID":2,"Name":"BrooksAguilar","Age":25,"About":"Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.\n","Gender":"male"},{"ID":27,"Name":"RebekahSutton","Age":26,"About":"Aliqua exercitation ad nostrud et exercitation amet quis cupidatat esse nostrud proident. Ullamco voluptate ex minim consectetur ea cupidatat in mollit reprehenderit voluptate labore sint laboris. Minim cillum et incididunt pariatur amet do esse. Amet irure elit deserunt quis culpa ut deserunt minim proident cupidatat nisi consequat ipsum.\n","Gender":"female"},{"ID":19,"Name":"BellBauer","Age":26,"About":"Nulla voluptate nostrud nostrud do ut tempor et quis non aliqua cillum in duis. Sit ipsum sit ut non proident exercitation. Quis consequat laboris deserunt adipisicing eiusmod non cillum magna.\n","Gender":"male"},{"ID":21,"Name":"JohnsWhitney","Age":26,"About":"Elit sunt exercitation incididunt est ea quis do ad magna. Commodo laboris nisi aliqua eu incididunt eu irure. Labore ullamco quis deserunt non cupidatat sint aute in incididunt deserunt elit velit. Duis est mollit veniam aliquip. Nulla sunt veniam anim et sint dolore.\n","Gender":"male"},{"ID":18,"Name":"TerrellHall","Age":27,"About":"Ut nostrud est est elit incididunt consequat sunt ut aliqua sunt sunt. Quis consectetur amet occaecat nostrud duis. Fugiat in irure consequat laborum ipsum tempor non deserunt laboris id ullamco cupidatat sit. Officia cupidatat aliqua veniam et ipsum labore eu do aliquip elit cillum. Labore culpa exercitation sint sint.\n","Gender":"male"},{"ID":20,"Name":"LoweryYork","Age":27,"About":"Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n","Gender":"male"},{"ID":3,"Name":"EverettDillard","Age":27,"About":"Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n","Gender":"male"},{"ID":8,"Name":"GlennJordan","Age":29,"About":"Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n","Gender":"male"},{"ID":5,"Name":"BeulahStark","Age":30,"About":"Enim cillum eu cillum velit labore. In sint esse nulla occaecat voluptate pariatur aliqua aliqua non officia nulla aliqua. Fugiat nostrud irure officia minim cupidatat laborum ad incididunt dolore. Fugiat nostrud eiusmod ex ea nulla commodo. Reprehenderit sint qui anim non ad id adipisicing qui officia Lorem.\n","Gender":"female"},{"ID":10,"Name":"HendersonMaxwell","Age":30,"About":"Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n","Gender":"male"},{"ID":4,"Name":"OwenLynn","Age":30,"About":"Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n","Gender":"male"},{"ID":22,"Name":"BethWynn","Age":31,"About":"Proident non nisi dolore id non. Aliquip ex anim cupidatat dolore amet veniam tempor non adipisicing. Aliqua adipisicing eu esse quis reprehenderit est irure cillum duis dolor ex. Laborum do aute commodo amet. Fugiat aute in excepteur ut aliqua sint fugiat do nostrud voluptate duis do deserunt. Elit esse ipsum duis ipsum.\n","Gender":"female"},{"ID":25,"Name":"KatherynJacobs","Age":32,"About":"Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n","Gender":"female"},{"ID":11,"Name":"GilmoreGuerra","Age":32,"About":"Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n","Gender":"male"},{"ID":28,"Name":"CohenHines","Age":32,"About":"Deserunt deserunt dolor ex pariatur dolore sunt labore minim deserunt. Tempor non et officia sint culpa quis consectetur pariatur elit sunt. Anim consequat velit exercitation eiusmod aute elit minim velit. Excepteur nulla excepteur duis eiusmod anim reprehenderit officia est ea aliqua nisi deserunt officia eiusmod. Officia enim adipisicing mollit et enim quis magna ea. Officia velit deserunt minim qui. Commodo culpa pariatur eu aliquip voluptate culpa ullamco sit minim laboris fugiat sit.\n","Gender":"male"},{"ID":30,"Name":"DicksonSilva","Age":32,"About":"Ipsum aliqua proident ullamco laboris eu occaecat deserunt. Amet ut adipisicing sint veniam dolore aliquip est mollit ex officia esse eiusmod veniam. Dolore magna minim aliquip sit deserunt. Nostrud occaecat dolore aliqua aliquip voluptate aliquip ad adipisicing.\n","Gender":"male"},{"ID":24,"Name":"GonzalezAnderson","Age":33,"About":"Quis consequat incididunt in ex deserunt minim aliqua ea duis. Culpa nisi excepteur sint est fugiat cupidatat nulla magna do id dolore laboris. Aute cillum eiusmod do amet dolore labore commodo do pariatur sit id. Do irure eiusmod reprehenderit non in duis sunt ex. Labore commodo labore pariatur ex minim qui sit elit.\n","Gender":"male"},{"ID":7,"Name":"LeannTravis","Age":34,"About":"Lorem magna dolore et velit ut officia. Cupidatat deserunt elit mollit amet nulla voluptate sit. Quis aute aliquip officia deserunt sint sint nisi. Laboris sit et ea dolore consequat laboris non. Consequat do enim excepteur qui mollit consectetur eiusmod laborum ut duis mollit dolor est. Excepteur amet duis enim laborum aliqua nulla ea minim.\n","Gender":"female"},{"ID":34,"Name":"KaneSharp","Age":34,"About":"Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n","Gender":"male"},{"ID":29,"Name":"ClarissaHenry","Age":34,"About":"Nostrud enim ea ad reprehenderit tempor ullamco exercitation. Elit in voluptate pariatur sit nisi occaecat laboris esse ipsum. Mollit elit et deserunt ea laboris sunt est amet culpa laboris occaecat ipsum sunt sunt.\n","Gender":"female"},{"ID":16,"Name":"AnnieOsborn","Age":35,"About":"Consequat fugiat veniam commodo nisi nostrud culpa pariatur. Aliquip velit adipisicing dolor et nostrud. Eu nostrud officia velit eiusmod ullamco duis eiusmod ad non do quis.\n","Gender":"female"},{"ID":12,"Name":"CruzGuerrero","Age":36,"About":"Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n","Gender":"male"},{"ID":9,"Name":"RoseCarney","Age":36,"About":"Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n","Gender":"female"},{"ID":17,"Name":"DillardMccoy","Age":36,"About":"Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n","Gender":"male"},{"ID":33,"Name":"TwilaSnow","Age":36,"About":"Sint non sunt adipisicing sit laborum cillum magna nisi exercitation. Dolore officia esse dolore officia ea adipisicing amet ea nostrud elit cupidatat laboris. Proident culpa ullamco aute incididunt aute. Laboris et nulla incididunt consequat pariatur enim dolor incididunt adipisicing enim fugiat tempor ullamco. Amet est ullamco officia consectetur cupidatat non sunt laborum nisi in ex. Quis labore quis ipsum est nisi ex officia reprehenderit ad adipisicing fugiat. Labore fugiat ea dolore exercitation sint duis aliqua.\n","Gender":"female"},{"ID":31,"Name":"PalmerScott","Age":37,"About":"Elit fugiat commodo laborum quis eu consequat. In velit magna sit fugiat non proident ipsum tempor eu. Consectetur exercitation labore eiusmod occaecat adipisicing irure consequat fugiat ullamco aliquip nostrud anim irure enim. Duis do amet cillum eiusmod eu sunt. Minim minim sunt sit sit enim velit sint tempor enim sint aliquip voluptate reprehenderit officia. Voluptate magna sit consequat adipisicing ut eu qui.\n","Gender":"male"},{"ID":26,"Name":"SimsCotton","Age":39,"About":"Ex cupidatat est velit consequat ad. Tempor non cillum labore non voluptate. Et proident culpa labore deserunt ut aliquip commodo laborum nostrud. Anim minim occaecat est est minim.\n","Gender":"male"},{"ID":6,"Name":"JenningsMays","Age":39,"About":"Veniam consectetur non non aliquip exercitation quis qui. Aliquip duis ut ad commodo consequat ipsum cupidatat id anim voluptate deserunt enim laboris. Sunt nostrud voluptate do est tempor esse anim pariatur. Ea do amet Lorem in mollit ipsum irure Lorem exercitation. Exercitation deserunt adipisicing nulla aute ex amet sint tempor incididunt magna. Quis et consectetur dolor nulla reprehenderit culpa laboris voluptate ut mollit. Qui ipsum nisi ullamco sit exercitation nisi magna fugiat anim consectetur officia.\n","Gender":"male"},{"ID":32,"Name":"ChristyKnapp","Age":40,"About":"Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n","Gender":"female"},{"ID":13,"Name":"WhitleyDavidson","Age":40,"About":"Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n","Gender":"male"}]`,
			StatusCode:      http.StatusOK,
		},

		// Проверка ошибок
		{
			TestDescribtion: "Invalid AccesssToken check",
			AccessToken:     "WrongToken",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"}]}",
			OrderField:      "Age",
			OrderBy:         1,
			Response:        "",
			StatusCode:      http.StatusUnauthorized,
		},
		{
			TestDescribtion: "Invalid OrderBy check",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"}]}",
			OrderField:      "Age",
			OrderBy:         12,
			Response:        `{"Error":"OrderBy invalid"}`,
			StatusCode:      http.StatusBadRequest,
		},
		{
			TestDescribtion: "Invalid OrderField filename check",
			AccessToken:     "hello",
			URL:             "http://localhost:8080/",
			Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"}]}",
			OrderField:      "WrongOrderField",
			OrderBy:         1,
			Response:        `{"Error":"OrderField invalid"}`,
			StatusCode:      http.StatusBadRequest,
		},
	}
	for caseNum, item := range cases {

		requestBody := SearchRequest{
			Limit:      25,
			Offset:     0,
			Query:      item.Query,
			OrderField: item.OrderField,
			OrderBy:    item.OrderBy,
		}

		body, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Println(fmt.Errorf("%s", err))
		}

		requestURL := item.URL + "?order_field=" + item.OrderField + "&query=" + item.Query + "&order_by=" + strconv.Itoa(item.OrderBy)
		req := httptest.NewRequest("GET", requestURL, bytes.NewReader(body))
		req.Header.Add("Accesstoken", item.AccessToken)

		w := httptest.NewRecorder()

		SearchServer(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d](%s) wrong StatusCode: got %d, expected %d",
				caseNum, item.TestDescribtion, w.Code, item.StatusCode)
		}

		resp := w.Result()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("[%d] failed to read body: %v", caseNum, err)
		}

		bodyStr := string(responseBody)

		if bodyStr != item.Response {
			t.Errorf("[%d](%s) wrong Response: got %+v, \nexpected %+v",
				caseNum, item.TestDescribtion, bodyStr, item.Response)
		}
	}
}

func TestSearchServerWrongFilename(t *testing.T) {

	caseFilename := ServerTestCase{
		TestDescribtion: "Invalid blank filename check",
		AccessToken:     "hello",
		URL:             "http://localhost:8080/",
		Query:           "{\"querylist\":[{\"Name\":\"GlennJordan\"}]}",
		OrderField:      "Age",
		OrderBy:         1,
		Response:        "",
		StatusCode:      http.StatusInternalServerError,
	}

	requestBody := SearchRequest{
		Limit:      25,
		Offset:     0,
		Query:      caseFilename.Query,
		OrderField: caseFilename.OrderField,
		OrderBy:    caseFilename.OrderBy,
	}

	filename = "WrongFileName"

	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println(fmt.Errorf("%s", err))
	}

	requestURL := caseFilename.URL + "?order_field=" + caseFilename.OrderField + "&query=" + caseFilename.Query + "&order_by=" + strconv.Itoa(caseFilename.OrderBy)
	req := httptest.NewRequest("GET", requestURL, bytes.NewReader(body))
	req.Header.Add("Accesstoken", caseFilename.AccessToken)

	w := httptest.NewRecorder()

	SearchServer(w, req)

	if w.Code != caseFilename.StatusCode {
		t.Errorf("(%s) wrong StatusCode: got %d, expected %d",
			caseFilename.TestDescribtion, w.Code, caseFilename.StatusCode)
	}

	resp := w.Result()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("failed to read body: %v", err)
	}

	bodyStr := string(responseBody)

	if bodyStr != caseFilename.Response {
		t.Errorf("(%s) wrong Response: got %+v, \nexpected %+v",
			caseFilename.TestDescribtion, bodyStr, caseFilename.Response)
	}

	filename = correctFilename
}

type ClientTestCase struct {
	TestDescribtion string
	Request         SearchRequest
	AccessToken     string
	Response        *SearchResponse
	IsError         bool
}

func TestFindUsers(t *testing.T) {
	cases := []ClientTestCase{
		{
			TestDescribtion: "fine request",
			IsError:         false,
			Request: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "Age",
				OrderBy:    1,
			},
			AccessToken: "hello",
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     8,
						Name:   "GlennJordan",
						Age:    29,
						About:  "Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n",
						Gender: "male",
					},
					{
						ID:     4,
						Name:   "OwenLynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
				},
			},
		},
		{
			TestDescribtion: "negative limit",
			IsError:         true,
			Request: SearchRequest{
				Limit:      -2,
				Offset:     0,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "Age",
				OrderBy:    1,
			},
			AccessToken: "hello",
			Response:    nil,
		},
		{
			TestDescribtion: "limit > 25",
			IsError:         false,
			Request: SearchRequest{
				Limit:      30,
				Offset:     0,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "Age",
				OrderBy:    1,
			},
			AccessToken: "hello",
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     8,
						Name:   "GlennJordan",
						Age:    29,
						About:  "Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n",
						Gender: "male",
					},
					{
						ID:     4,
						Name:   "OwenLynn",
						Age:    30,
						About:  "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n",
						Gender: "male",
					},
				},
			},
		},
		{
			TestDescribtion: "negative offset",
			IsError:         true,
			Request: SearchRequest{
				Limit:      2,
				Offset:     -2,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "Age",
				OrderBy:    1,
			},
			AccessToken: "hello",
			Response:    nil,
		},
		{
			TestDescribtion: "Bad AccessToken",
			IsError:         true,
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "Age",
				OrderBy:    1,
			},
			AccessToken: "BadAccessToken",
			Response:    nil,
		},
		{
			TestDescribtion: "Bad OrderField",
			IsError:         true,
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "BadOrderField",
				OrderBy:    1,
			},
			AccessToken: "hello",
			Response:    nil,
		},
		{
			TestDescribtion: "NextPage == true",
			IsError:         false,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "Age",
				OrderBy:    1,
			},
			AccessToken: "hello",
			Response: &SearchResponse{
				Users: []User{
					{
						ID:     8,
						Name:   "GlennJordan",
						Age:    29,
						About:  "Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
		},
		{
			TestDescribtion: "Incorrect orderBy",
			IsError:         true,
			Request: SearchRequest{
				Limit:      15,
				Offset:     0,
				Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
				OrderField: "Age",
				OrderBy:    111,
			},
			AccessToken: "hello",
			Response:    nil,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for caseNum, item := range cases {

		client := &SearchClient{
			AccessToken: item.AccessToken,
			URL:         ts.URL,
		}

		result, err := client.FindUsers(item.Request)

		// fmt.Println(caseNum, "RESULT:", result)

		if err != nil && !item.IsError {
			t.Errorf("[%d](%s) unexpected error: %#v", caseNum, item.TestDescribtion, err)
		}
		if err == nil && item.IsError {
			t.Errorf("[%d](%s) expected error, got nil", caseNum, item.TestDescribtion)
		}
		// if
		if !reflect.DeepEqual(item.Response, result) {
			t.Errorf("[%d](%s) wrong result, EXPECTED: %#v, GOT: %#v", caseNum, item.TestDescribtion, item.Response, result)
		}
	}
	ts.Close()
}

func TestFindUnknownClientDo(t *testing.T) {
	myCase := ClientTestCase{
		TestDescribtion: "Unknown error in client.Do()",
		IsError:         true,
		Request: SearchRequest{
			Limit:      3,
			Offset:     0,
			Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField: "Age",
			OrderBy:    1,
		},
		AccessToken: "hello",
		Response:    nil,
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		AccessToken: myCase.AccessToken,
		URL:         "2128506",
	}
	result, err := client.FindUsers(myCase.Request)

	if err != nil && !myCase.IsError {
		t.Errorf("(%s) unexpected error: %#v", myCase.TestDescribtion, err)
	}
	if err == nil && myCase.IsError {
		t.Errorf("(%s) expected error, got nil", myCase.TestDescribtion)
	}
	if !reflect.DeepEqual(myCase.Response, result) {
		t.Errorf("(%s) wrong result, EXPECTED: %#v, GOT: %#v", myCase.TestDescribtion, myCase.Response, result)
	}
	ts.Close()
}

func TestFindTimeout(t *testing.T) {
	myCase := ClientTestCase{
		TestDescribtion: "Timeout",
		IsError:         true,
		Request: SearchRequest{
			Limit:      3,
			Offset:     0,
			Query:      "{\"querylist\":[{\"Name\":\"нучтоеще\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField: "Age",
			OrderBy:    1,
		},
		AccessToken: "hello",
		Response:    nil,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second)
	}))

	client := &SearchClient{
		AccessToken: myCase.AccessToken,
		URL:         ts.URL,
	}
	result, err := client.FindUsers(myCase.Request)

	if err != nil && !myCase.IsError {
		t.Errorf("(%s) unexpected error: %#v", myCase.TestDescribtion, err)
	}
	if err == nil && myCase.IsError {
		t.Errorf("(%s) expected error, got nil", myCase.TestDescribtion)
	}
	if !reflect.DeepEqual(myCase.Response, result) {
		t.Errorf("(%s) wrong result, EXPECTED: %#v, GOT: %#v", myCase.TestDescribtion, myCase.Response, result)
	}
	ts.Close()
}

func TestFindUsersFilename(t *testing.T) {
	myCase := ClientTestCase{
		TestDescribtion: "Bad filename in server",
		IsError:         true,
		Request: SearchRequest{
			Limit:      3,
			Offset:     0,
			Query:      "{\"querylist\":[{\"Name\":\"GlennJordan\"},{\"Name\":\"OwenLynn\"}]}",
			OrderField: "Age",
			OrderBy:    1,
		},
		AccessToken: "hello",
		Response:    nil,
	}

	filename = "InvalidFilename"
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		AccessToken: myCase.AccessToken,
		URL:         ts.URL,
	}
	result, err := client.FindUsers(myCase.Request)

	// fmt.Println(caseNum, "RESULT:", result)

	if err != nil && !myCase.IsError {
		t.Errorf("(%s) unexpected error: %#v", myCase.TestDescribtion, err)
	}
	if err == nil && myCase.IsError {
		t.Errorf("(%s) expected error, got nil", myCase.TestDescribtion)
	}
	// if
	if !reflect.DeepEqual(myCase.Response, result) {
		t.Errorf("(%s) wrong result, EXPECTED: %#v, GOT: %#v", myCase.TestDescribtion, myCase.Response, result)
	}
	ts.Close()
}
