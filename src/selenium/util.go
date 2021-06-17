package selenium

import (
	"fmt"
	"sifamaGO/src/db"
	"sifamaGO/src/model"
	"sifamaGO/src/util"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

func ProperTitle(input string) string {
	words := strings.Fields(input)

	for index, word := range words {
		if len(word) < 4 {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

func KeepMouseMoving(quit chan string) {
	for {
		select {
		case <-quit:
			return
		default:
			robotgo.MoveMouse(100, 300)
			time.Sleep(time.Minute * 3)
			robotgo.MoveMouse(300, 500)
			time.Sleep(time.Second * 3)
		}
	}

}

func IsTrechosDNIT(km float64) bool {

	return (km >= 211.29 && km <= 230.06) || (km >= 277 && km <= 360)
}

func IsLocationValid(caption string, local *model.Local) (string, error) {

	kmInicial := local.KmInicialDouble
	kmFinal := local.KmFinalDouble
	palavraChave := local.Tro.PalavraChave
	oldKmInicial := local.KmInicial
	oldKmFinal := local.KmFinal
	var checkkmInicial bool
	var checkkmFinal bool
	var dnit bool
	var dnit1 bool

	if strings.Contains(local.Rodovia, "364") && strings.Contains(strings.ToLower(local.Sentido), "decrescente") {
		if (kmInicial > 0 && kmInicial < 20) || (kmFinal > 0 && kmFinal < 20) {
			newKmInicial, newKmFinal := InterpolationLocal(local)
			kmFinalStr := fmt.Sprintf("%.3f", newKmFinal)
			kmInicialStr := fmt.Sprintf("%.3f", newKmInicial)
			kmFinalStr = strings.Replace(kmFinalStr, ".", ",", -1)
			kmInicialStr = strings.Replace(kmInicialStr, ".", ",", -1)

			local.KmFinal = kmFinalStr
			local.KmFinalDouble = newKmFinal
			local.KmInicial = kmInicialStr
			local.KmInicialDouble = newKmInicial
			caption = caption + " (km da 364 Variante : " + oldKmInicial + " - " + oldKmFinal + " )"
			db.GetDB().Save(&local)
		}

	}

	kmInicial = local.KmInicialDouble
	kmFinal = local.KmFinalDouble

	//***

	checkkmInicial, dnit, err = CheckKm(local.Estado, local.Rodovia, kmInicial, palavraChave)
	if err != nil {
		return "", err
	}
	checkkmFinal, dnit1, err = CheckKm(local.Estado, local.Rodovia, kmFinal, palavraChave)
	if err != nil {
		return "", err
	}

	local.Valid = checkkmInicial && checkkmFinal
	local.TrechoDNIT = dnit || dnit1

	db.GetDB().Save(&local)

	return caption, nil
}

func InterpolationLocal(local *model.Local) (float64, float64) {

	kmInicial := local.KmInicialDouble
	kmFinal := local.KmFinalDouble
	newKmInicial := InterpolationKm(kmInicial)
	newKmFinal := InterpolationKm(kmFinal)
	if newKmInicial < newKmFinal {
		newKmInicial, newKmFinal = newKmFinal, newKmInicial
	}
	return newKmInicial, newKmFinal
}
func InterpolationKm(km float64) float64 {
	if km <= 10.94 {
		return 360 - ((km - 1.0) / 9.94 * 7.7)
	} else if km <= 12.25 {
		return 351 + (12.25 - km)
	} else if km <= 19.67 {
		return 351 - ((km-12.25)/7.42)*7.2
	} else {
		return km
	}

}

func CheckKm(estado, rodovia string, km float64, palavraChave string) (bool, bool, error) {
	if util.CONCESSIONARIA == "CRO" {
		return CheckKmCRO(rodovia, km, palavraChave)
	} else if util.CONCESSIONARIA == "MSVIA" {
		return CheckKMSVIA(rodovia, km, palavraChave)
	} else if util.CONCESSIONARIA == "ECO050" {
		return CheckKmECO050(estado, rodovia, km, palavraChave)
	}
	return false, false, fmt.Errorf("nao foi possivel verificar o Check-km .. ")
}
func CheckKMSVIA(rodovia string, km float64, palavraChave string) (bool, bool, error) {
	if km >= 0 && km <= 847.2 {
		return true, false, nil
	}
	return false, false, nil
}

func CheckKmECO050(estado, rodovia string, km float64, palavraChave string) (bool, bool, error) {
	switch rodovia {
	case "50":
		if estado == "MG" {
			if (km >= 0 && km <= 65.7) || (km >= 77.3 && km <= 207.3) {
				return true, false, nil
			}
		} else if estado == "GO" {
			if km >= 95.7 && km <= 314.2 {
				return true, false, nil
			}
		}

	case "Contorno de Uberlândia":
		return (km >= 0 && km <= 22.4), false, nil

	default:
		return false, false, nil
	}
	return false, false, nil
}

func CheckKmCRO(rodovia string, km float64, palavraChave string) (bool, bool, error) {

	disposicao := GetDisposicaoLegal(palavraChave)
	edificacoes := false

	if disposicao[0] == "5" && disposicao[1] == "742" {
		edificacoes = true
	}

	switch rodovia {
	case "70":
		if km < 495.9 || km > 524 {
			return false, false, nil
		}
		return true, false, nil

	case "163":
		return (km < 120 || (km <= 855 && km >= 507.1)), false, nil

	case "364":
		if km < 201 || km > 588.2 {
			return false, false, nil
		} else if km <= 211.3 {
			return true, false, nil
		} else if km >= 434.6 {
			return true, false, nil
		} else if !edificacoes && IsTrechosDNIT(km) {

			return true, true, nil
		} else if edificacoes && (km >= 201 && km <= 588.2) {
			return true, false, nil
		}

	default:
		return false, false, nil
	}
	return false, false, nil
}

func GetDisposicaoLegal(palavraChave string) []string {

	palavraChave = strings.ToLower(palavraChave)

	if strings.Contains(palavraChave, "buraco") {
		return []string{"6", "774"}
	} else if strings.Contains(palavraChave, "afundamen") || strings.Contains(palavraChave, "escorregamento") || strings.Contains(palavraChave, "remendo") {
		return []string{"6", "773"}
	} else if strings.Contains(palavraChave, "drenagem") {
		return []string{"6", "782"}
	} else if strings.Contains(palavraChave, "meio fio") || strings.Contains(palavraChave, "meio-fio") {
		return []string{"4", "720"}
	} else if strings.Contains(palavraChave, "desplacam") || strings.Contains(palavraChave, "deformaç") {
		return []string{"6", "773"}
	} else if strings.Contains(palavraChave, "vertical") || strings.Contains(palavraChave, "horizontal") {
		return []string{"7", "807"}
	} else if strings.Contains(palavraChave, "terrapleno") || strings.Contains(palavraChave, "talude") {
		return []string{"6", "783"}
	} else if strings.Contains(palavraChave, "defensa") {
		return []string{"7", "808"}
	} else if strings.Contains(palavraChave, "instalac") || strings.Contains(palavraChave, "instalaç") || strings.Contains(palavraChave, "edifica") {
		return []string{"5", "742"}
	} else if strings.Contains(palavraChave, "pmv") {
		return []string{"5", "752"}
	} else if strings.Contains(palavraChave, "guarda corpo") || strings.Contains(palavraChave, "guarda-corpo") {
		return []string{"7", "811"}
	} else if strings.Contains(palavraChave, "sujeira") {
		return []string{"7", "806"}
	}

	disposicaoLegal := make(map[string][]string)

	disposicaoLegal["4-v"] = []string{"4", "718"}
	disposicaoLegal["4-vi"] = []string{"4", "719"}
	disposicaoLegal["4-vii"] = []string{"4", "720"}
	disposicaoLegal["4-xii"] = []string{"4", "725"}
	disposicaoLegal["4-xiii"] = []string{"4", "726"}

	disposicaoLegal["5-iii"] = []string{"5", "742"}
	disposicaoLegal["5-v"] = []string{"5", "744"}
	disposicaoLegal["5-ix"] = []string{"5", "748"}
	disposicaoLegal["5-xii"] = []string{"5", "751"}
	disposicaoLegal["5-xiii"] = []string{"5", "752"}
	disposicaoLegal["5-xiv"] = []string{"5", "753"}
	disposicaoLegal["5-xv"] = []string{"5", "754"}
	disposicaoLegal["5-xxviii"] = []string{"5", "767"}

	disposicaoLegal["6-iii"] = []string{"6", "773"}
	disposicaoLegal["6-iv"] = []string{"6", "774"}
	disposicaoLegal["6-v"] = []string{"6", "775"}
	disposicaoLegal["6-vii"] = []string{"6", "777"}
	disposicaoLegal["6-viii"] = []string{"6", "778"}
	disposicaoLegal["6-x"] = []string{"6", "780"}
	disposicaoLegal["6-xi"] = []string{"6", "781"}
	disposicaoLegal["6-xii"] = []string{"6", "782"}
	disposicaoLegal["6-xiii"] = []string{"6", "783"}
	disposicaoLegal["6-xiv"] = []string{"6", "784"}
	disposicaoLegal["6-xvi"] = []string{"6", "786"}
	disposicaoLegal["6-xvii"] = []string{"6", "787"}
	disposicaoLegal["6-xxviii"] = []string{"6", "798"}

	disposicaoLegal["7-viii"] = []string{"7", "806"}
	disposicaoLegal["7-ix"] = []string{"7", "807"}
	disposicaoLegal["7-x"] = []string{"7", "808"}
	disposicaoLegal["7-xii"] = []string{"7", "810"}
	disposicaoLegal["7-xiii"] = []string{"7", "811"}

	disposicaoLegal["8-vii"] = []string{"8", "864"}

	disposicaoLegal["9-vii"] = []string{"9", "863"}

	return disposicaoLegal[palavraChave]
}

func GetDescricaoDisposicaoLegal(codstr string) string {

	cod, err := strconv.Atoi(codstr)
	errorHandle(err)

	descricao := make(map[int]string)

	descricao[718] = "Art. 4º, V - deixar selagem em juntas de pavimento rígido ou trincas em desconformidade com o PER, por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER"
	descricao[719] = "Art. 4º, VI - deixar de manter marcos quilométricos ou mantê-los em más condições de visibilidade, por prazo superior a 7 (sete) dias, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER"
	descricao[720] = "Art. 4º, VII - deixar meios-fios danificados, deteriorados ou ausentes por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER"
	descricao[725] = "Art. 4º, XII - deixar barreira de concreto de Obra-de-Arte Especial - OAE sem pintura por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER"
	descricao[726] = "Art. 4º, XIII - deixar armaduras de OAE sem recobrimento por prazo superior a 48 (quarenta e oito horas)"

	// Art 5

	descricao[742] = "Art 5º, III - deixar de executar os serviços de conservação das instalações, áreas operacionais e bens vinculados à concessão por prazo superior a 72 horas após a ocorrência de evento que comprometa suas condições normais de uso e a integridade do bem"
	descricao[744] = "Art 5º, V - deixar de remover, da faixa de domínio, material resultante de poda, capina ou obras no prazo de 48 (quarenta e oito) horas, salvo no caso de materiais reaproveitáveis ou de bota-foras autorizados pela ANTT"
	descricao[748] = "Art 5º, IX - deixar de repor ou manter tachas, tachões e balizadores refletivos danificados ou ausentes no prazo de 72 (setenta e duas) horas"
	descricao[751] = "Art 5º, XII - deixar de adotar medidas, ainda que provisórias, para reparação de cercamento nas áreas operacionais por prazo superior a 24 (vinte e quatro) horas"
	descricao[752] = "Art 5º, XIII - deixar de adotar medidas, ainda que provisórias, para reparar painel de mensagem variável inoperante ou em condições que não permitam a transmissão de informações aos usuários, por prazo superior a 72 (setenta e duas) horas"
	descricao[752] = "Art 5º, XIV - deixar de adotar medidas, ainda que provisórias para reparação das cercas limítrofes da faixa de proteção e de seus aceiros por prazo superior a 72 (setenta e duas) horas"
	descricao[754] = "Art 5º, XV - deixar de adotar medidas, ainda que provisórias, para corrigir falha em sistema ou equipamento dos postos de pesagem no prazo de 24 (vinte e quatro) horas ou de acordo com o especificado no Contrato e/ou PER, se este fizer referência diversa"
	descricao[767] = "Art 5º, XXVIII - deixar de adotar providências para corrigir desnível entre faixas contíguas, ainda que em caráter provisório, no prazo de 24 (vinte e quatro) horas, ou, deixar de implementar a solução definitiva para correção no prazo estabelecido pela ANTT"

	// Art 6:

	descricao[773] = "Art 6º, III - deixar de corrigir depressões, abaulamentos (escorregamentos de massa asfáltica) ou áreas exsudadas na pista ou no acostamento, no prazo de 72 (setenta e duas) horas, ou conforme previsto no Contrato de Concessão e/ou PER"
	descricao[774] = "Art 6º, IV - deixar de corrigir/tapar buracos, panelas na pista ou no acostamento, no prazo de 24 (vinte e quatro) horas, ou conforme previsto no Contrato de Concessão e/ou PER"
	descricao[775] = "Art 6º, V - deixar de corrigir, no pavimento rígido, defeitos com grau de severidade alto, no prazo de 7 (sete) dias, ou conforme previsto no Contrato de Concessão e/ou PER"
	descricao[777] = "Art 6º, VII - deixar de corrigir, no pavimento rígido, defeitos de alçamento de placa, fissura de canto, placa dividida (rompida), escalonamento ou degrau, placa bailarina, quebras localizadas e buracos no prazo de 48 (quarenta e oito) horas, ou conforme previsto no Contrato de Concessão e/ou PER"
	descricao[778] = "Art 6º, VIII - deixar de manter ou manter de forma não visível pelos usuários sinalização (vertical ou aérea) de indicação, de serviços auxiliares ou educativas, por prazo superior a 7 (sete) dias"
	descricao[780] = "Art 6º, X - deixar de manter ou manter de forma não funcional dispositivo anti-ofuscante por prazo superior a 7 (sete) dias, ou conforme previsto no Contrato de Concessão ou no PER"
	descricao[781] = "Art 6º, XI - deixar com problemas de conservação elemento de OAE, exceto guarda-corpo, por prazo superior a 30 (trinta) dias ou conforme Contrato de Concessão e/ou PER"
	descricao[782] = "Art 6º, XII - deixar de reparar, limpar ou desobstruir sistema de drenagem e Obra-de-Arte Corrente-OAC por prazo superior a 72 (setenta e duas) horas, ou conforme previsto no Contrato de Concessão ou no PER"
	descricao[783] = "Art 6º, XIII - deixar de adotar providências para solucionar, ainda que de modo provisório, processo erosivo ou condição de instabilidade em talude, por prazo superior a 72 (setenta e duas) horas, ou deixar de implementar solução definitiva no prazo estabelecido pela ANTT"
	descricao[784] = "Art 6º, XIV - deixar de manter ou manter de forma não funcional o sistema de iluminação da rodovia, por prazo superior a 48 (quarenta e oito) horas"
	descricao[786] = "Art 6º, XVI - deixar de corrigir falha em equipamento de praça de pedágio no prazo de 6 (seis) horas, sem prejuízo ao atendimento dos parâmetros de desempenho estabelecidos no PER"
	descricao[787] = `Art 6º, XVII - deixar "Call Box" inoperante por prazo superior a 24 (vinte e quatro) horas, ou de acordo com o especificado no PER, se este fizer referência diversa`
	descricao[798] = "Art 6º, XXVIII - deixar de intervir, mesmo que provisoriamente, em recalque em pavimento na cabeceira de OAE e/ou OAC por prazo superior a 72 (setenta e duas) horas, desde que essa obrigação tenha sido prevista no Contrato de Concessão ou PER"

	// Art 7

	descricao[806] = "Art 7º, VIII - deixar de remover material da(s) faixa(s) de rolamento( s) ou acostamento(s) que obstrua ou comprometa a correta fluidez do tráfego no prazo de 6 (seis) horas a partir do evento que lhe deu origem"
	descricao[807] = "Art 7º, IX - deixar de manter ou manter a sinalização horizontal, vertical ou aérea, em desconformidade com as normas técnicas vigentes, por prazo superior ao estabelecido pela ANTT, excluídas as ocorrências previstas nos artigos 5°, 6° e 9°"
	descricao[808] = "Art 7º, X - deixar de recompor barreira rígida ou defensa metálica danificada no prazo de 48 horas"
	descricao[810] = "Art 7º, 	XII - deixar de intervir para restaurar a funcionalidade de elemento da rodovia quando da ocorrência de fatos oriundos da ação de terceiros ou de eventos da natureza que possam colocar em risco a segurança do usuário, no prazo de 48 (quarenta e oito) horas ou conforme estabelecido pela ANTT"
	descricao[811] = "Art 7º, XIII - deixar de recuperar, ainda que provisoriamente, guarda- corpo de OAE, inclusive passarela, por prazo superior a 24 (vinte e quatro) horas, ou, deixar de efetuar sua reposição definitiva, por prazo superior a 72 (setenta e duas) horas, ou conforme Contrato e/ou PER"

	// Art 8

	descricao[864] = "Art 8º, VII - deixar de adotar as providências cabíveis, inclusive por vias judiciais, para garantia do patrimônio da rodovia, da faixa de domínio, das edificações e dos bens da concessão, inclusive quanto à implantação de acessos irregulares e ocupações ilegais; Nos casos de constatação destas irregularidades para as concessões da 2ª etapa, há previsão contratual de prazo de 24 (vinte e quatro) horas para a correção. Deste modo, deverá ser expedido TRO enquadrado neste mesmo Art. 8º, inciso VII, da Resolução"

	// Art 9

	descricao[863] = "Art 9º, VII - deixar de manter ou manter sinalização vertical de regulamentação em desconformidade com as normas técnicas vigentes, por prazo superior ao previsto no Contrato de Concessão ou no PER"

	return descricao[cod]
}

func getEstadoERodovia(word string, i int) (string, string, string, error) {

	if strings.Contains(word, "mt") && (strings.Contains(word, "364") || strings.Contains(word, "70") || strings.Contains(word, "163")) {
		util.CONCESSIONARIA = "CRO"
		return getEstadoERodoviaCRO(word, i)
	} else if strings.Contains(word, "ms") && strings.Contains(word, "163") {
		util.CONCESSIONARIA = "MSVIA"
		return getEstadoERodoviaMsVia(word, i)
	} else if strings.Contains(word, "torno") || strings.Contains(word, "50") && (strings.Contains(word, "go") || strings.Contains(word, "mg")) {
		util.CONCESSIONARIA = "ECO050"
		return getEstadoERodoviaEco050(word, i)
	} else {
		return "", "", "", fmt.Errorf("nao foi possivel identificar a concessionária a partir do primeiro local")
	}

}

func CheckKmEco050(estado, rodovia string, km float64, palavraChave string) (bool, bool) {

	switch rodovia {
	case "50":
		if estado == "MG" {
			if (km >= 0 && km <= 65.7) || (km >= 77.3 && km <= 207.3) {
				return true, false
			}
		} else if estado == "GO" {
			if km >= 95.7 && km <= 314.2 {
				return true, false
			}
		}

	case "Contorno de Uberlândia":
		return (km >= 0 && km <= 22.4), false

	default:
		return false, false
	}
	return false, false
}

func getEstadoERodoviaCRO(word string, i int) (string, string, string, error) {
	var rodovia string
	estado := "MT"
	concessionaria := "CRO"
	if strings.Contains(word, "070") {
		rodovia = "70"
		return concessionaria, rodovia, estado, nil
	} else if strings.Contains(word, "163") {
		rodovia = "163"
		return concessionaria, rodovia, estado, nil
	} else if strings.Contains(word, "364") {
		rodovia = "364"
		return concessionaria, rodovia, estado, nil
	} else {
		err := fmt.Errorf("erro. não consegui identificar Rodovia na linha %v.... abortando", i+1)
		return concessionaria, rodovia, estado, err
	}

}

func getEstadoERodoviaEco050(word string, i int) (string, string, string, error) {
	var estado string
	var rodovia string
	concessionaria := "ECO050"

	if strings.Contains(word, "050/mg") && !strings.Contains(word, "torno") {
		rodovia = "50"
		estado = "MG"
		return concessionaria, rodovia, estado, nil

	} else if strings.Contains(word, "050/go") {
		rodovia = "50"
		estado = "GO"
		return concessionaria, rodovia, estado, nil
	} else if strings.Contains(word, "contorno") {
		rodovia = "Contorno de Uberlândia"
		estado = "MG"
		return concessionaria, rodovia, estado, nil
	} else {
		err := fmt.Errorf("erro. não consegui identificar Rodovia na linha %v.... abortando", i+1)
		return concessionaria, rodovia, estado, err
	}
}
func getEstadoERodoviaMsVia(word string, i int) (string, string, string, error) {
	var estado string
	var rodovia string
	concessionaria := "MSVIA"
	if strings.Contains(word, "163") {
		rodovia = "163"
		estado = "MS"
		return concessionaria, rodovia, estado, nil
	}
	err := fmt.Errorf("erro. não consegui identificar Rodovia na linha %v.... abortando", i+1)
	return concessionaria, rodovia, estado, err
}

//art 4
//        718	V - deixar selagem em juntas de pavimento rígido ou trincas em desconformidade com o PER, por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER
//        719	VI - deixar de manter marcos quilométricos ou mantê-los em más condições de visibilidade, por prazo superior a 7 (sete) dias, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER
//        720	VII - deixar meios-fios danificados, deteriorados ou ausentes por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER
//        725	XII - deixar barreira de concreto de Obra-de-Arte Especial - OAE sem pintura por prazo superior a 72 (setenta e duas) horas, ou conforme prazo diverso previsto no Contrato de Concessão ou no PER
//        726	XIII - deixar armaduras de OAE sem recobrimento por prazo superior a 48 (quarenta e oito horas)

// Art 5

//      742	III - deixar de executar os serviços de conservação das instalações, áreas operacionais e bens vinculados à concessão por prazo superior a 72 horas após a ocorrência de evento que comprometa suas condições normais de uso e a integridade do bem
//      744	V - deixar de remover, da faixa de domínio, material resultante de poda, capina ou obras no prazo de 48 (quarenta e oito) horas, salvo no caso de materiais reaproveitáveis ou de bota-foras autorizados pela ANTT
//      748	IX - deixar de repor ou manter tachas, tachões e balizadores refletivos danificados ou ausentes no prazo de 72 (setenta e duas) horas
//      751	XII - deixar de adotar medidas, ainda que provisórias, para reparação de cercamento nas áreas operacionais por prazo superior a 24 (vinte e quatro) horas
//      752	XIII - deixar de adotar medidas, ainda que provisórias, para reparar painel de mensagem variável inoperante ou em condições que não permitam a transmissão de informações aos usuários, por prazo superior a 72 (setenta e duas) horas
//      753	XIV - deixar de adotar medidas, ainda que provisórias para reparação das cercas limítrofes da faixa de proteção e de seus aceiros por prazo superior a 72 (setenta e duas) horas
//      754	XV - deixar de adotar medidas, ainda que provisórias, para corrigir falha em sistema ou equipamento dos postos de pesagem no prazo de 24 (vinte e quatro) horas ou de acordo com o especificado no Contrato e/ou PER, se este fizer referência diversa
//      767	XXVIII - deixar de adotar providências para corrigir desnível entre faixas contíguas, ainda que em caráter provisório, no prazo de 24 (vinte e quatro) horas, ou, deixar de implementar a solução definitiva para correção no prazo estabelecido pela ANTT

//Art 6:

//       773	III - deixar de corrigir depressões, abaulamentos (escorregamentos de massa asfáltica) ou áreas exsudadas na pista ou no acostamento, no prazo de 72 (setenta e duas) horas, ou conforme previsto no Contrato de Concessão e/ou PER
//      774	IV - deixar de corrigir/tapar buracos, panelas na pista ou no acostamento, no prazo de 24 (vinte e quatro) horas, ou conforme previsto no Contrato de Concessão e/ou PER
//      775	V - deixar de corrigir, no pavimento rígido, defeitos com grau de severidade alto, no prazo de 7 (sete) dias, ou conforme previsto no Contrato de Concessão e/ou PER
//      777	VII - deixar de corrigir, no pavimento rígido, defeitos de alçamento de placa, fissura de canto, placa dividida (rompida), escalonamento ou degrau, placa bailarina, quebras localizadas e buracos no prazo de 48 (quarenta e oito) horas, ou conforme previsto no Contrato de Concessão e/ou PER
//      778	VIII - deixar de manter ou manter de forma não visível pelos usuários sinalização (vertical ou aérea) de indicação, de serviços auxiliares ou educativas, por prazo superior a 7 (sete) dias
//      780	X - deixar de manter ou manter de forma não funcional dispositivo anti-ofuscante por prazo superior a 7 (sete) dias, ou conforme previsto no Contrato de Concessão ou no PER
//      781	XI - deixar com problemas de conservação elemento de OAE, exceto guarda-corpo, por prazo superior a 30 (trinta) dias ou conforme Contrato de Concessão e/ou PER
//      782	XII - deixar de reparar, limpar ou desobstruir sistema de drenagem e Obra-de-Arte Corrente-OAC por prazo superior a 72 (setenta e duas) horas, ou conforme previsto no Contrato de Concessão ou no PER
//      783	XIII - deixar de adotar providências para solucionar, ainda que de modo provisório, processo erosivo ou condição de instabilidade em talude, por prazo superior a 72 (setenta e duas) horas, ou deixar de implementar solução definitiva no prazo estabelecido pela ANTT
//      784	XIV - deixar de manter ou manter de forma não funcional o sistema de iluminação da rodovia, por prazo superior a 48 (quarenta e oito) horas
//      786	XVI - deixar de corrigir falha em equipamento de praça de pedágio no prazo de 6 (seis) horas, sem prejuízo ao atendimento dos parâmetros de desempenho estabelecidos no PER
//      787	XVII - deixar "Call Box" inoperante por prazo superior a 24 (vinte e quatro) horas, ou de acordo com o especificado no PER, se este fizer referência diversa
//      798	XXVIII - deixar de intervir, mesmo que provisoriamente, em recalque em pavimento na cabeceira de OAE e/ou OAC por prazo superior a 72 (setenta e duas) horas, desde que essa obrigação tenha sido prevista no Contrato de Concessão ou PER

//Art 7

//        806	VIII - deixar de remover material da(s) faixa(s) de rolamento( s) ou acostamento(s) que obstrua ou comprometa a correta fluidez do tráfego no prazo de 6 (seis) horas a partir do evento que lhe deu origem
//          807	IX - deixar de manter ou manter a sinalização horizontal, vertical ou aérea, em desconformidade com as normas técnicas vigentes, por prazo superior ao estabelecido pela ANTT, excluídas as ocorrências previstas nos artigos 5°, 6° e 9°
//        808	X - deixar de recompor barreira rígida ou defensa metálica danificada no prazo de 48 horas
//          810	XII - deixar de intervir para restaurar a funcionalidade de elemento da rodovia quando da ocorrência de fatos oriundos da ação de terceiros ou de eventos da natureza que possam colocar em risco a segurança do usuário, no prazo de 48 (quarenta e oito) horas ou conforme estabelecido pela ANTT
//          811	XIII - deixar de recuperar, ainda que provisoriamente, guarda- corpo de OAE, inclusive passarela, por prazo superior a 24 (vinte e quatro) horas, ou, deixar de efetuar sua reposição definitiva, por prazo superior a 72 (setenta e duas) horas, ou conforme Contrato e/ou PER

// Art 8

//        864	VII - deixar de adotar as providências cabíveis, inclusive por vias judiciais, para garantia do patrimônio da rodovia, da faixa de domínio, das edificações e dos bens da concessão, inclusive quanto à implantação de acessos irregulares e ocupações ilegais; Nos casos de constatação destas irregularidades para as concessões da 2ª etapa, há previsão contratual de prazo de 24 (vinte e quatro) horas para a correção. Deste modo, deverá ser expedido TRO enquadrado neste mesmo Art. 8º, inciso VII, da Re</option>

// Art 9

//        863	VII - deixar de manter ou manter sinalização vertical de regulamentação em desconformidade com as normas técnicas vigentes, por prazo superior ao previsto no Contrato de Concessão ou no PER
