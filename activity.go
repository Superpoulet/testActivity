package sample

import (
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

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

	output := &Output{AnOutput: input.AnInput}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
