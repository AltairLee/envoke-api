package spec

import (
	. "github.com/zbo14/envoke/common"
	regex "github.com/zbo14/envoke/regex"
)

const CONTEXT = "http://localhost:8888/spec#Context"

func NewLink(id string) Data {
	return Data{"@id": id}
}

func GetId(data Data) string {
	return data.GetStr("@id")
}

func SetId(data Data, id string) {
	data.Set("@id", id)
}

func MatchId(id string) bool {
	return MatchStr(regex.ID, id)
}

func GetType(data Data) string {
	return data.GetStr("@type")
}

func NewParty(email, ipi, isni string, memberIds []string, name, pro, sameAs, _type string) Data {
	party := Data{
		"@context": CONTEXT,
		"@type":    _type,
		"name":     name,
	}
	switch _type {
	case "MusicGroup", "Organization":
		if n := len(memberIds); n > 0 {
			member := make([]Data, n)
			for i, memberId := range memberIds {
				if !MatchId(memberId) {
					panic("Invalid memberId")
				}
				member[i] = NewLink(memberId)
			}
			party.Set("member", member)
		}
	case "Person":
		//..
	default:
		panic(ErrorAppend(ErrInvalidType, _type))
	}
	if MatchStr(regex.EMAIL, email) {
		party.Set("email", email)
	}
	if MatchStr(regex.IPI, ipi) {
		party.Set("ipiNumber", ipi)
	}
	if MatchStr(regex.ISNI, isni) {
		party.Set("isniNumber", isni)
	}
	if MatchStr(regex.PRO, pro) {
		party.Set("pro", pro)
	}
	if MatchUrlRelaxed(sameAs) {
		party.Set("sameAs", sameAs)
	}
	return party
}

func GetDescription(data Data) string {
	return data.GetStr("description")
}

func GetEmail(data Data) string {
	return data.GetStr("email")
}

func GetIPI(data Data) string {
	return data.GetStr("ipiNumber")
}

func GetISNI(data Data) string {
	return data.GetStr("isniNumber")
}

func GetName(data Data) string {
	return data.GetStr("name")
}

func GetPRO(data Data) string {
	return data.GetStr("pro")
}

func GetSameAs(data Data) string {
	return data.GetStr("sameAs")
}

func NewCollaboration(memberIds []string, name string, roleNames, splits []string) Data {
	if len(memberIds) == 0 {
		panic("No memberIds")
	}
	var haveRoleNames bool
	if n := len(roleNames); n > 0 {
		if n != len(memberIds) {
			panic("Number of roles doesn't equal number of memberIds")
		}
		haveRoleNames = true
	}
	var haveSplits bool
	if n := len(splits); n > 0 {
		if n != len(splits) {
			panic("Number of splits doesn't equal number of memberIds")
		}
		haveSplits = true
	}
	member := make([]Data, len(memberIds))
	for i, memberId := range memberIds {
		if !MatchId(memberId) {
			panic("Invalid memberId")
		}
		member[i] = Data{
			"@type":  "OrganizationRole",
			"member": NewLink(memberId),
		}
		if haveRoleNames {
			member[i].Set("roleName", roleNames[i])
		}
		if haveSplits {
			member[i].Set("split", splits[i])
		}
	}
	collaboration := Data{
		"@context": CONTEXT,
		"@type":    "MusicCollaboration",
		"member":   member,
	}
	if !EmptyStr(name) {
		collaboration.Set("name", name)
	}
	return collaboration
}

func GetCollaboratorIds(data Data) []string {
	members := data.GetDataSlice("member")
	memberIds := make([]string, len(members))
	for i, member := range members {
		memberIds[i] = GetId(member.GetData("member"))
	}
	return memberIds
}

func NewComposition(composerId, hfa, iswc, lang, name, publisherId, sameAs string) Data {
	composition := Data{
		"@context": CONTEXT,
		"@type":    "MusicComposition",
		"composer": NewLink(composerId),
		"name":     name,
	}
	if MatchStr(regex.HFA, hfa) {
		composition.Set("hfaCode", hfa)
	}
	if MatchStr(regex.ISWC, iswc) {
		composition.Set("iswcCode", iswc)
	}
	if MatchStr(regex.LANGUAGE, lang) {
		composition.Set("inLanguage", lang)
	}
	if MatchId(publisherId) {
		composition.Set("publisher", NewLink(publisherId))
	}
	if MatchUrlRelaxed(sameAs) {
		composition.Set("sameAs", sameAs)
	}
	return composition
}

func GetComposerId(data Data) string {
	composer := data.GetData("composer")
	return GetId(composer)
}

func GetHFA(data Data) string {
	return data.GetStr("hfaCode")
}

func GetISWC(data Data) string {
	return data.GetStr("iswcCode")
}

func GetLanguage(data Data) string {
	return data.GetStr("inLanguage")
}

func GetPublisherId(data Data) string {
	publisher := data.GetData("publisher")
	return GetId(publisher)
}

func NewPublication(compositionIds []string, compositionRightIds []string, name, publisherId, sameAs string) Data {
	m := len(compositionIds)
	if m == 0 {
		panic("No compositionIds")
	}
	compositions := make([]Data, m)
	for i, compositionId := range compositionIds {
		compositions[i] = Data{
			"@type":    "ListItem",
			"position": i + 1,
			"item": Data{
				// "@type": "MusicComposition",
				"@id": compositionId,
			},
		}
	}
	n := len(compositionRightIds)
	if n == 0 {
		panic("No compositionRightIds")
	}
	compositionRights := make([]Data, n)
	for i, compositionRightId := range compositionRightIds {
		compositionRights[i] = Data{
			"@type":    "ListItem",
			"position": i + 1,
			"item": Data{
				// "@type": "CompositionRight",
				"@id": compositionRightId,
			},
		}
	}
	publication := Data{
		"@context": CONTEXT,
		"@type":    "MusicPublication",
		"composition": Data{
			"@type":           "ItemList",
			"numberOfItems":   m,
			"itemListElement": compositions,
		},
		"compositionRight": Data{
			"@type":           "ItemList",
			"numberOfItems":   n,
			"itemListElement": compositionRights,
		},
		"name":      name,
		"publisher": NewLink(publisherId),
	}
	if MatchUrlRelaxed(sameAs) {
		publication.Set("sameAs", sameAs)
	}
	return publication
}

func GetCompositionIds(data Data) []string {
	compositions := data.GetData("composition")
	n := compositions.GetInt("numberOfItems")
	compositionIds := make([]string, n)
	itemListElement := compositions.GetInterfaceSlice("itemListElement")
	for i, elem := range itemListElement {
		item := AssertData(elem).GetData("item")
		compositionIds[i] = GetId(item)
	}
	return compositionIds
}

func GetCompositionRightIds(data Data) []string {
	compositionRights := data.GetData("compositionRight")
	n := compositionRights.GetInt("numberOfItems")
	compositionRightIds := make([]string, n)
	itemListElement := compositionRights.GetInterfaceSlice("itemListElement")
	for i, elem := range itemListElement {
		item := AssertData(elem).GetData("item")
		compositionRightIds[i] = GetId(item)
	}
	return compositionRightIds
}

func NewRecording(artistId, compositionId, compositionRightId, duration, isrc, mechanicalLicenseId, publicationId, recordLabelId, sameAs string) Data {
	recording := Data{
		"@context":    CONTEXT,
		"@type":       "MusicRecording",
		"byArtist":    NewLink(artistId),
		"recordingOf": NewLink(compositionId),
	}
	if MatchId(compositionRightId) {
		if !MatchId(publicationId) {
			panic("must have compositionRightId and publicationId")
		}
		recording.Set("compositionRight", NewLink(compositionRightId))
		recording.Set("publication", NewLink(publicationId))
	} else if MatchId(mechanicalLicenseId) {
		recording.Set("mechanicalLicense", NewLink(mechanicalLicenseId))
	} else {
		// artist should be composer
	}
	if !EmptyStr(duration) {
		recording.Set("duration", duration)
	}
	if MatchStr(regex.ISRC, isrc) {
		recording.Set("isrcCode", isrc)
	}
	if MatchId(recordLabelId) {
		recording.Set("recordLabel", NewLink(recordLabelId))
	}
	if MatchUrlRelaxed(sameAs) {
		recording.Set("sameAs", sameAs)
	}
	return recording
}

func GetArtistId(data Data) string {
	artist := data.GetData("byArtist")
	return GetId(artist)
}

func GetCompositionRightId(data Data) string {
	compositionRight := data.GetData("compositionRight")
	return GetId(compositionRight)
}

func GetMechanicalLicenseId(data Data) string {
	mechanicalLicense := data.GetData("mechanicalLicense")
	return GetId(mechanicalLicense)
}

func GetProducerId(data Data) string {
	producer := data.GetData("producer")
	return GetId(producer)
}

func GetPublicationId(data Data) string {
	publication := data.GetData("publication")
	return GetId(publication)
}

func GetRecordingOfId(data Data) string {
	composition := data.GetData("recordingOf")
	return GetId(composition)
}

func GetRecordLabelId(data Data) string {
	recordLabel := data.GetData("recordLabel")
	return GetId(recordLabel)
}

func NewRelease(name string, recordingIds, recordingRightIds []string, recordLabelId, sameAs string) Data {
	m := len(recordingIds)
	if m == 0 {
		panic("No recordingIds")
	}
	recordings := make([]Data, m)
	for i, recordingId := range recordingIds {
		recordings[i] = Data{
			"@type":    "ListItem",
			"position": i + 1,
			"item": Data{
				// "@type": "MusicRecording",
				"@id": recordingId,
			},
		}
	}
	n := len(recordingRightIds)
	if n == 0 {
		panic("No recordingRightIds")
	}
	recordingRights := make([]Data, n)
	for i, recordingRightId := range recordingRightIds {
		recordingRights[i] = Data{
			"@type":    "ListItem",
			"position": i + 1,
			"item": Data{
				// "@type": "RecordingRight",
				"@id": recordingRightId,
			},
		}
	}
	release := Data{
		"@context": CONTEXT,
		"@type":    "MusicRelease",
		"name":     name,
		"recording": Data{
			"@type":           "ItemList",
			"numberOfItems":   m,
			"itemListElement": recordings,
		},
		"recordingRight": Data{
			"@type":           "ItemList",
			"numberOfItems":   n,
			"itemListElement": recordingRights,
		},
		"recordLabel": NewLink(recordLabelId),
	}
	if MatchUrlRelaxed(sameAs) {
		release.Set("sameAs", sameAs)
	}
	return release
}

func GetRecordingIds(data Data) []string {
	recordings := data.GetData("recording")
	n := recordings.GetInt("numberOfItems")
	recordingIds := make([]string, n)
	itemListElement := recordings.GetInterfaceSlice("itemListElement")
	for i, elem := range itemListElement {
		item := AssertData(elem).GetData("item")
		recordingIds[i] = GetId(item)
	}
	return recordingIds
}

func GetRecordingRightIds(data Data) []string {
	recordingRights := data.GetData("recordingRight")
	n := recordingRights.GetInt("numberOfItems")
	recordingRightIds := make([]string, n)
	itemListElement := recordingRights.GetInterfaceSlice("itemListElement")
	for i, elem := range itemListElement {
		item := AssertData(elem).GetData("item")
		recordingRightIds[i] = GetId(item)
	}
	return recordingRightIds
}

// Note: percentageShares is taken from the tx output amount so it's not included in the data model

func NewCompositionRight(recipientId, senderId string, territory []string, validFrom, validThrough string) Data {
	return NewRight(recipientId, senderId, territory, "CompositionRight", validFrom, validThrough)
}

func NewRecordingRight(recipientId, senderId string, territory []string, validFrom, validThrough string) Data {
	return NewRight(recipientId, senderId, territory, "RecordingRight", validFrom, validThrough)
}

func NewRight(recipientId, senderId string, territory []string, _type, validFrom, validThrough string) Data {
	return Data{
		"@context":     CONTEXT,
		"@type":        _type,
		"recipient":    NewLink(recipientId),
		"sender":       NewLink(senderId),
		"territory":    territory,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
}

func GetRecipientId(data Data) string {
	recipient := data.GetData("recipient")
	return GetId(recipient)
}

func GetRecipientShares(data Data) int {
	return data.GetInt("recipientShares")
}

func GetSenderId(data Data) string {
	sender := data.GetData("sender")
	return GetId(sender)
}

func GetSenderShares(data Data) int {
	return data.GetInt("senderShares")
}

func GetTerritory(data Data) []string {
	return data.GetStrSlice("territory")
}

// Note: txId is the hex id of a TRANSFER tx in Bigchain/IPDB
// the output amount(s) will specify shares transferred/kept

func NewCompositionRightTransfer(compositionRightId, publicationId, recipientId, senderId, txId string) Data {
	return Data{
		"@context":         CONTEXT,
		"@type":            "CompositionRightTransfer",
		"compositionRight": NewLink(compositionRightId),
		"publication":      NewLink(publicationId),
		"recipient":        NewLink(recipientId),
		"sender":           NewLink(senderId),
		"tx":               NewLink(txId),
	}
}

func GetCompositionRightTransferId(data Data) string {
	compositionRightTransfer := data.GetData("compositionRightTransfer")
	return GetId(compositionRightTransfer)
}

func GetTxId(data Data) string {
	tx := data.GetData("tx")
	return GetId(tx)
}

func NewRecordingRightTransfer(recipientId, recordingRightId, releaseId, senderId, txId string) Data {
	return Data{
		"@context":       CONTEXT,
		"@type":          "RecordingRightTransfer",
		"recipient":      NewLink(recipientId),
		"recordingRight": NewLink(recordingRightId),
		"release":        NewLink(releaseId),
		"sender":         NewLink(senderId),
		"tx":             NewLink(txId),
	}
}

func GetReleaseId(data Data) string {
	release := data.GetData("release")
	return GetId(release)
}

func GetRecordingRightTransferId(data Data) string {
	recordingRightTransfer := data.GetData("recordingRightTransfer")
	return GetId(recordingRightTransfer)
}

func NewMechanicalLicense(compositionIds []string, compositionRightId, compositionRightTransferId, publicationId, recipientId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	mechanicalLicense := Data{
		"@context":     CONTEXT,
		"@type":        "MechanicalLicense",
		"recipient":    NewLink(recipientId),
		"sender":       NewLink(senderId),
		"territory":    territory,
		"usage":        usage,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
	n := len(compositionIds)
	if n > 0 {
		compositions := make([]Data, n)
		for i, compositionId := range compositionIds {
			if !MatchId(compositionId) {
				panic(ErrorAppend(ErrInvalidId, compositionId))
			}
			compositions[i] = Data{
				"@type":    "ListItem",
				"position": i + 1,
				"item": Data{
					// "@type": "MusicComposition",
					"@id": compositionId,
				},
			}
		}
		mechanicalLicense.Set("composition", Data{
			"@type":           "ItemList",
			"numberOfItems":   n,
			"itemListElement": compositions,
		})
	} else if !MatchId(publicationId) {
		panic("Expected valid compositionIds or publicationId")
	}
	if MatchId(publicationId) {
		mechanicalLicense.Set("publication", NewLink(publicationId))
		if MatchId(compositionRightId) {
			mechanicalLicense.Set("compositionRight", NewLink(compositionRightId))
		} else if MatchId(compositionRightTransferId) {
			mechanicalLicense.Set("compositionRightTransfer", NewLink(compositionRightTransferId))
		} else {
			panic("Expected valid compositionRightId or compositionRightTransferId")
		}
	}
	return mechanicalLicense
}

func NewMasterLicense(recipientId string, recordingIds []string, recordingRightId, recordingRightTransferId, releaseId, senderId string, territory, usage []string, validFrom, validThrough string) Data {
	masterLicense := Data{
		"@context":     CONTEXT,
		"@type":        "MasterLicense",
		"recipient":    NewLink(recipientId),
		"sender":       NewLink(senderId),
		"territory":    territory,
		"usage":        usage,
		"validFrom":    validFrom,
		"validThrough": validThrough,
	}
	n := len(recordingIds)
	if n > 0 {
		recordings := make([]Data, n)
		for i, recordingId := range recordingIds {
			recordings[i] = Data{
				"@type":    "ListItem",
				"position": i + 1,
				"item": Data{
					// "@type": "MusicRecording",
					"@id": recordingId,
				},
			}
		}
		masterLicense.Set("recording", Data{
			"@type":           "ItemList",
			"numberOfItems":   n,
			"itemListElement": recordings,
		})
	} else if !MatchId(releaseId) {
		panic("Expected valid recordingIds or releaseId")
	}
	if MatchId(releaseId) {
		masterLicense.Set("release", NewLink(releaseId))
		if MatchId(recordingRightId) {
			masterLicense.Set("recordingRight", NewLink(recordingRightId))
		} else if MatchId(recordingRightTransferId) {
			masterLicense.Set("recordingRightTransfer", NewLink(recordingRightTransferId))
		} else {
			panic("Expected valid recordingRightId or recordingRightTransferId")
		}
	}
	return masterLicense
}

func GetRecordingRightId(data Data) string {
	recordingRight := data.GetData("recordingRight")
	return GetId(recordingRight)
}
