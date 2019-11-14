package sample

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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

	/*	ctx.Logger().Debugf("Input: %s", input.IdContrat)*/
	log.Info("idContrat : ", input)
	log.Info("DateDebutRecherche : ", input.DateDebutRecherche)
	log.Info("DateFinRecherche : ", input.DateFinRecherche)

	critere := CriteresRecherchePaiement{

		DateDebutRecherche: input.DateDebutRecherche,
		DateFinRecherche:   input.DateFinRecherche,
		NumeroContrat:      input.IdContrat,
	}
	log.Info("criteres recherche : ", critere)
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
	log.Info("SOAP Call: lirePaimentPrestaAI")
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
	resultJson := result.SoapBody.LirePaiementPrestationOut

	resJson, err := json.Marshal(resultJson)
	if err != nil {
		log.Info("Cannot encode to JSON ", err)
	}
	/*	buftes := new(bytes.Buffer)
		buftes.ReadFrom(res.Body)
		newStr := buftes.String()
	*/
	log.Debug("SOAP result json: ", resJson)

	output := &Output{string(resJson)}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
