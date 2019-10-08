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
