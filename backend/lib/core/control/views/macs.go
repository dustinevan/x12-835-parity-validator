package views

import (
	"context"

	"github.com/dustinevan/x12-835-parity-validator/backend/pkg/macs"
)

// The MAC workload numbers as a tree (/api/macs): backend/pkg/macs Tree,
// with JSON names. Static reference data, like Codes: it reads no table and
// does not depend on which submissions are loaded.

// MACsResponse is every jurisdiction in CMS's order: the A/B MACs, then the
// DME MACs.
type MACsResponse struct {
	Jurisdictions []MACJurisdictionView `json:"jurisdictions"`
}

// MACJurisdictionView is one MAC jurisdiction, its contractor, and its
// Parts: A then B, or DME alone.
type MACJurisdictionView struct {
	Jurisdiction string        `json:"jurisdiction"`
	Contractor   string        `json:"contractor"`
	Parts        []MACPartView `json:"parts"`
}

// MACPartView is one claim Part ("A", "B", "DME", or "" when the table
// gives none): its roll-ups with the workloads under each, the workloads
// under no roll-up (a DME jurisdiction's number is one of these), the home
// health and hospice numbers, and the numbers the numbering does not place.
type MACPartView struct {
	Part      string          `json:"part"`
	RollUps   []MACRollUpView `json:"roll_ups"`
	Workloads []MACNumberView `json:"workloads"`
	HHH       []MACNumberView `json:"hhh"`
	Other     []MACNumberView `json:"other"`
}

// MACRollUpView is a roll-up number and the workloads under it.
// CarriesHHH is true when the CR says the roll-up also carries the home
// health and hospice claims.
type MACRollUpView struct {
	MACNumberView
	CarriesHHH bool            `json:"carries_hhh"`
	Workloads  []MACNumberView `json:"workloads"`
}

// MACNumberView is one five-digit number, the states it covers, and a note
// from the table's comments, or null.
type MACNumberView struct {
	ID     string         `json:"id"`
	States []MACStateView `json:"states"`
	Note   nullable       `json:"note"`
}

// MACStateView is a USPS code and its name.
type MACStateView struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// MACs returns backend/pkg/macs Tree. The tree is built on each call, which
// the one view that serves it does rarely.
func (v Controller) MACs(ctx context.Context) (MACsResponse, error) {
	tree := macs.Tree()
	out := MACsResponse{Jurisdictions: make([]MACJurisdictionView, 0, len(tree))}
	for _, j := range tree {
		jv := MACJurisdictionView{
			Jurisdiction: j.Name,
			Contractor:   j.Contractor,
			Parts:        make([]MACPartView, 0, len(j.Parts)),
		}
		for _, p := range j.Parts {
			pv := MACPartView{
				Part:      p.Name,
				RollUps:   make([]MACRollUpView, 0, len(p.RollUps)),
				Workloads: macNumberViews(p.Workloads),
				HHH:       macNumberViews(p.HHH),
				Other:     macNumberViews(p.Other),
			}
			for _, r := range p.RollUps {
				pv.RollUps = append(pv.RollUps, MACRollUpView{
					MACNumberView: macNumberView(r.Number),
					CarriesHHH:    r.CarriesHHH,
					Workloads:     macNumberViews(r.Workloads),
				})
			}
			jv.Parts = append(jv.Parts, pv)
		}
		out.Jurisdictions = append(out.Jurisdictions, jv)
	}
	return out, nil
}

// macNumberViews converts a list of numbers; an empty list is [] in JSON,
// never null.
func macNumberViews(ns []macs.Number) []MACNumberView {
	out := make([]MACNumberView, 0, len(ns))
	for _, n := range ns {
		out = append(out, macNumberView(n))
	}
	return out
}

func macNumberView(n macs.Number) MACNumberView {
	states := make([]MACStateView, 0, len(n.States))
	for _, s := range n.States {
		states = append(states, MACStateView{Code: s.Code, Name: s.Name})
	}
	return MACNumberView{ID: n.ID, States: states, Note: nullable(n.Note)}
}
