package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Rule struct {
	value       string
	isCondition bool
	left, right *Rule
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need arg")
	}
	contentsBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("%s", err)
	}
	contents := string(contentsBytes)
	rules := make(map[string]*Rule)
	parts := make([]map[string]int, 0)
	isRules := true
	for _, line := range strings.Split(contents, "\n") {
		if len(line) == 0 {
			isRules = false
			continue
		}
		if isRules {
			lineParts := strings.Split(line, "{")
			ruleKey := lineParts[0]
			ruleValue := lineParts[1]
			ruleValue = ruleValue[:len(ruleValue)-1]
			rules[ruleKey] = parseRules(ruleValue)
		} else {
			part := make(map[string]int)
			line = line[1 : len(line)-1]
			lineParts := strings.Split(line, ",")
			for _, item := range lineParts {
				itemParts := strings.Split(item, "=")
				k := itemParts[0]
				v, _ := strconv.Atoi(itemParts[1])
				part[k] = v
			}
			parts = append(parts, part)
		}
	}

	result := 0
	for i, part := range parts {
		rule := rules["in"]
		r := rule.evaluate(part)
		for r != "A" && r != "R" {
			rule = rules[r]
			r = rule.evaluate(part)
		}
		if r == "A" {
			log.Printf("%d: accepted", i)
			for _, v := range part {
				result += v
			}
		} else {
			log.Printf("%d: rejected", i)
		}
	}
	log.Printf("result: %d", result)
}

func parseRules(rules string) (res *Rule) {
	ruleStack := make([]*Rule, 0)
	i := 0
	j := 0
	for j < len(rules) {
		if rules[j] == ':' {
			ruleStack[len(ruleStack)-1].value = rules[i:j]
			i = j + 1
		}
		if rules[j] == ',' {
			ruleStack = append(ruleStack, &Rule{value: rules[i:j]})
			i = j + 1
		}
		if rules[j] == '<' || rules[j] == '>' {
			rule := &Rule{isCondition: true}
			ruleStack = append(ruleStack, rule)
			ruleStack[len(ruleStack)-1].isCondition = true
		}
		j++
	}
	if i < len(rules) {
		ruleStack = append(ruleStack, &Rule{value: rules[i:]})
	} else {
		log.Printf("did not append tail of rules: %s", rules)
	}
	for len(ruleStack) > 1 {
		for i := 0; i < len(ruleStack); i++ {
			if ruleStack[i].isCondition && !ruleStack[i].isReady() {
				if i+2 >= len(ruleStack) {
					log.Printf("overflow")
					for k, rule := range ruleStack {
						log.Printf("%d: %s", k, rule.value)
					}
					log.Fatalf(rules)
				}
				if ruleStack[i+1].isReady() && ruleStack[i+2].isReady() {
					ruleStack[i].left = ruleStack[i+1]
					ruleStack[i].right = ruleStack[i+2]
					if i+3 < len(ruleStack) {
						copy(ruleStack[i+1:], ruleStack[i+3:])
					}
					ruleStack = ruleStack[:len(ruleStack)-2]
				}
			}
		}
	}
	res = ruleStack[0]
	return
}

func (r *Rule) evaluate(part map[string]int) string {
	//log.Printf("eval rule: %s", r.value)
	for r.isCondition {
		isLess := true
		vParts := strings.Split(r.value, "<")
		if len(vParts) == 1 {
			isLess = false
			vParts = strings.Split(r.value, ">")
		}
		left := 0
		right := 0
		c, err := strconv.Atoi(vParts[0])
		l := vParts[1]
		if err == nil {
			left = c
			lv := part[l]
			right = lv
		} else {
			c, _ = strconv.Atoi(vParts[1])
			right = c
			l = vParts[0]
			lv := part[l]
			left = lv
		}
		if isLess {
			if left < right {
				r = r.left
			} else {
				r = r.right
			}
		} else {
			if left > right {
				r = r.left
			} else {
				r = r.right
			}
		}
	}
	return r.value
}

func (r *Rule) isReady() bool {
	if !r.isCondition {
		return true
	}
	if r.isCondition && r.left != nil && r.right != nil {
		return true
	}
	return false
}
