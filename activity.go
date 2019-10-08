package sample

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"

	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

var log = logger.GetLogger("activity-Soap-call")

func init() {
	_ = activity.Register(&Activity{}) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Debugf("Setting: %s", s.ASetting)

	act := &Activity{} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
}

type CriteresRecherchePaiement struct {
	DateDebutRecherche string `xml:"dateDebutRecherche,omitempty"`

	DateFinRecherche string `xml:"dateFinRecherche,omitempty"`

	NumeroContrat string `xml:"numeroContrat,omitempty"`
}

type CodeSystemeExterne string

type LirePaiementPrestationIn struct {
	CodeSystemeExterne        CodeSystemeExterne        `xml:"codeSystemeExterne"`
	CriteresRecherchePaiement CriteresRecherchePaiement `xml:"criteresRecherchePaiement,omitempty"`
}

type SoapBody struct {
	LirePaiementPrestationIn LirePaiementPrestationIn `xml:"pai:lirePaiementPrestation"`
}

type SoapEnveloppe struct {
	XMLName  xml.Name `xml:"soapenv:Envelope"`
	XmlNS    string   `xml:"xmlns:soapenv,attr"`
	XmlNS2   string   `xml:"xmlns:pai,attr"`
	SoapBody SoapBody `xml:"soapenv:Body"`
}

type SoapEnvelopeOUT struct {
	XMLName  xml.Name    `xml:"Envelope"`
	SoapBody SoapBodyOut `xml:"Body"`
}
type SoapBodyOut struct {
	LirePaiementPrestationOut LirePaiementPrestationOut `xml:"lirePaiementPrestationOut"`
}
type LirePaiementPrestationOut struct {
	XMLName       xml.Name      `xml:"http://wsi.cegedimactiv.com/client/paiementPrestation lirePaiementPrestationOut"`
	ListePaiement ListePaiement `xml:"listePaiement,omitempty"`
}

type ListePaiement struct {
	NombreDePaiements int32 `xml:"nombreDePaiements,omitempty"`

	PaiementsParContrat []PaiementsParContrat `xml:"PaiementsParContrat,omitempty"`
}

type PaiementsParContrat struct {
	DatePaiement string `xml:"datePaiement,omitempty"`

	Destinataire *Destinataire `xml:"Destinataire,omitempty"`

	LibelleModePaiement string `xml:"libelleModePaiement,omitempty"`

	ModePaiement string `xml:"modePaiement,omitempty"`

	MontantPaiement float32 `xml:"montantPaiement,omitempty"`

	NumeroPaiement int `xml:"numeroPaiement,omitempty"`
}

type Destinataire struct {
	Civilite string `xml:"civilite,omitempty"`

	Nom string `xml:"nom,omitempty"`

	Prenom string `xml:"prenom,omitempty"`
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	ctx.Logger().Debugf("Input: %s", input.AnInput)

	critere := CriteresRecherchePaiement{

		DateDebutRecherche: "2017-01-01",
		DateFinRecherche:   "2018-01-01",
		NumeroContrat:      "10140540",
	}

	request := LirePaiementPrestationIn{
		CodeSystemeExterne:        "UNEO",
		CriteresRecherchePaiement: critere,
	}
	Body := SoapBody{
		LirePaiementPrestationIn: request,
	}

	Envelope := SoapEnveloppe{
		XmlNS:    "http://schemas.xmlsoap.org/soap/envelope/",
		XmlNS2:   "http://wsi.cegedimactiv.com/client/paiementPrestation",
		SoapBody: Body,
	}

	var buf bytes.Buffer

	enc := xml.NewEncoder(&buf)
	enc.Indent("  ", "    ")

	if err := enc.Encode(Envelope); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	url := fmt.Sprintf("%s%s",
		"https://inf1-unur04.priv.services-fm.net",
		"/e-services-client/services/paiementPrestation",
	)

	soapAction := "urn:lirePaiementPrestation" // The format is `urn:<soap_action>`
	httpMethod := "POST"

	req, err := http.NewRequest(httpMethod, url, &buf)
	if err != nil {
		log.Debug("Error ", err.Error())
		return
	}
	req.Header.Set("Content-type", "text/xml;charset=UTF-8")
	req.Header.Set("SOAPAction", soapAction)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	log.Debug("SOAP Call: lirePaimentPrestaAI")
	res, err := client.Do(req)
	if err != nil {
		log.Debug("Error on dispatching request. ", err.Error())
		return
	}
	result := new(SoapEnvelopeOUT)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Debug("Error on unmarshaling xml. ", err.Error())
		return
	}
	log.Debug("SOAP result: ", result)
	/*	bodyString := string(body)
		fmt.Println(bodyString)*/

	output := &Output{LirePaiementPrestationOut: input.AnInput}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
