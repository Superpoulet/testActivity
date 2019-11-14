package sample

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	ASetting string `md:"aSetting,required"`
}

type Input struct {
	IdContrat          string `md:"idContrat,required"`
	DateDebutRecherche string `md:"dateDebutRecherche,required"`
	DateFinRecherche   string `md:"dateFinRecherche,required"`
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
	/*XMLName       xml.Name      `xml:"lirePaiementPrestationOut"`*/
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

func (r *Input) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["idContrat"])
	strVal2, _ := coerce.ToString(values["dateDebutRecherche"])
	strVal3, _ := coerce.ToString(values["dateFinRecherche"])
	r.IdContrat = strVal
	r.DateDebutRecherche = strVal2
	r.DateFinRecherche = strVal3
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"idContrat":          r.IdContrat,
		"dateDebutRecherche": r.DateDebutRecherche,
		"dateFinRecherche":   r.DateFinRecherche,
	}
}

type Output struct {
	LirePaiementPrestationOut string `md:"LirePaiementPrestationOut"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["result"])
	o.LirePaiementPrestationOut = strVal
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"LirePaiementPrestationOut": o.LirePaiementPrestationOut,
	}
}
