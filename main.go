package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
    "github.com/hyperjumptech/grule-rule-engine/pkg"
) 

type Entity struct {
	CustomerID		uint 
	SaleValue		float32
	ConcurrentSale	bool
	SaleDate		time.Time
    Office			string
    CumulatedPoints  float32
}

func main() {

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request){
       entity := &Entity{
           CustomerID:1,
           SaleValue:1000.2,
           ConcurrentSale:true,
           SaleDate: time.Now(),
           Office:"London",
       }

       dataCtx := ast.NewDataContext()
        err := dataCtx.Add("MF", entity)
        if err != nil {
            panic(err)
        }
       knowledgeBase := loadRules() 
       eng1 := &engine.GruleEngine{MaxCycle: 1}
       err = eng1.Execute(dataCtx, knowledgeBase)
       if err != nil {
           panic(err)
       }
        fmt.Fprintf(w,"This sale cumulate you %f points", entity.CumulatedPoints)
    })

	fmt.Printf("Starting server at port 80\n")
	if err := http.ListenAndServe(":80", nil); err != nil {
        log.Fatal(err)
    }
}


func loadRules() *ast.KnowledgeBase {
    knowledgeLibrary := ast.NewKnowledgeLibrary()
    ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
    // lets prepare a rule definition
    drls := `
    rule LondonOffice "London Office" salience 10{
        when
            MF.SaleValue > 500 && MF.Office == "London" && MF.ConcurrentSale
        then
            MF.CumulatedPoints = (MF.SaleValue*10/1000);
            Retract("LondonOffice");
    }

    rule NYOfficeNew "New York Office" salience 10{
        when
            MF.SaleValue > 500 && MF.Office == "New York" && MF.ConcurrentSale
        then
            MF.CumulatedPoints = (MF.SaleValue*10/1000);
            Retract("NYOfficeNew");
    }

    rule NYOffice "New York Office 2" salience 9{
        when
            MF.SaleValue > 500 && MF.Office == "New York" && MF.ConcurrentSale && MF.CumulatedPoints > 0
        then
            MF.CumulatedPoints = (MF.SaleValue*8/1000);
            Retract("NYOffice");
    }
   
            
    `
    // Add the rule definition above into the library and name it 'TutorialRules'  version '0.0.1'
    bs := pkg.NewBytesResource([]byte(drls))
    err := ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
    if err != nil {
        panic(err)
    }
    knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

    return knowledgeBase
}