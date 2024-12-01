package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Rule struct {
	value          string
	isCondition    bool
	left, right    *Rule
	isLess         bool
	isVariableLeft bool
	constant       int
	variable       string
}

type Part map[string]int

type Parts struct {
	parts []int
	label string
}

type Args struct {
	iMap map[string]int
	jMap map[string]int
	rule *Rule
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
	//parts := make([]Part, 0)
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
			/*
				part := make(Part)
				line = line[1 : len(line)-1]
				lineParts := strings.Split(line, ",")
				for _, item := range lineParts {
					itemParts := strings.Split(item, "=")
					k := itemParts[0]
					v, _ := strconv.Atoi(itemParts[1])
					part[k] = v
				}
				parts = append(parts, part)
			*/
		}
	}

	log.Printf("rules: %d", len(rules))
	labelToParts := make(map[string]*Parts)
	labels := []string{"x", "m", "a", "s"}
	for _, k := range labels {
		ps := make([]int, 0, 4000)
		for i := 1; i <= 4000; i++ {
			ps = append(ps, i)
		}
		psObj := &Parts{
			parts: ps,
			label: k,
		}
		labelToParts[k] = psObj
		log.Printf("created label %s with %d parts", k, len(ps))
	}

	result := 0
	rootRule := rules["in"]
	iMap := make(map[string]int)
	jMap := make(map[string]int)
	for _, l := range labels {
		iMap[l] = 0
		jMap[l] = 4000
	}
	args := &Args{
		iMap: iMap,
		jMap: jMap,
		rule: rootRule,
	}
	stack := make([]*Args, 0)
	stack = append(stack, args)
	for len(stack) > 0 {
		args := stack[len(stack)-1]
		//log.Printf("at rule: %s", args.rule.value)
		if !args.rule.isCondition {
			if args.rule.value == "A" {
				// TODO need to dedup
				r := 1
				for _, l := range labels {
					i := args.iMap[l]
					j := args.jMap[l]
					r *= j - i
				}
				//log.Printf("rule is accept: adding %d", r)
				result += r
				stack = stack[:len(stack)-1]
				continue
			}
			if args.rule.value == "R" {
				//log.Printf("rule is reject")
				stack = stack[:len(stack)-1]
				continue
			}
			//log.Printf("rule is a pointer")
			newRule := rules[args.rule.value]
			args.rule = newRule
			continue
		}
		//log.Printf("rule is a condition, evaluating")
		stack = stack[:len(stack)-1]
		variable := args.rule.variable
		ps := labelToParts[variable]
		partHi, ok := args.bisect(ps)
		if !ok {
			log.Fatal("bad bisect")
		}
		argsHi := args.makeCopy()
		argsHi.iMap[variable] = partHi
		argsHi.rule = args.rule.right
		argsLo := args.makeCopy()
		argsLo.jMap[variable] = partHi
		argsLo.rule = args.rule.left
		stack = append(stack, argsHi)
		stack = append(stack, argsLo)
	}
	log.Printf("%d", result)
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
				j := 0
				leftOperand := ""
				rightOperand := ""
				for k, c := range ruleStack[i].value {
					if c == '<' || c == '>' {
						if c == '<' {
							ruleStack[i].isLess = true
						} else {
							ruleStack[i].isLess = false
						}
						leftOperand = ruleStack[i].value[j:k]
						j = k + 1
					}
				}
				rightOperand = ruleStack[i].value[j:]
				v, err := strconv.Atoi(leftOperand)
				if err != nil {
					ruleStack[i].isVariableLeft = true
					v, _ = strconv.Atoi(rightOperand)
					ruleStack[i].variable = leftOperand
					ruleStack[i].constant = v
				} else {
					ruleStack[i].isVariableLeft = false
					ruleStack[i].variable = rightOperand
					ruleStack[i].constant = v
				}
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

/*
func (p Parts) Less(i, j int) bool {
	iv := p.parts[i][p.label]
	jv := p.parts[j][p.label]
	return iv < jv
}

func (p Parts) Swap(i, j int) {
	p.parts[i], p.parts[j] = p.parts[j], p.parts[i]
}

func (p Parts) Len() int {
	return len(p.parts)
}
*/

func (a *Args) bisect(p *Parts) (mid int, ok bool) {
	if !a.rule.isCondition || a.rule.variable != p.label {
		ok = false
		return
	}
	r := a.rule
	constant := r.constant
	lo := a.iMap[r.variable]
	hi := a.jMap[r.variable]
	if r.isLess == r.isVariableLeft {
		for lo < hi {
			mid = (lo + hi) / 2
			midv := p.parts[mid]
			if midv <= constant {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
	} else if r.isLess != r.isVariableLeft {
		for lo < hi {
			mid = (lo + hi) / 2
			midv := p.parts[mid]
			if midv > r.constant {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
	}
	ok = true
	return
}

func (a *Args) makeCopy() (b *Args) {
	iMap := make(map[string]int)
	jMap := make(map[string]int)
	for k := range a.iMap {
		iMap[k] = a.iMap[k]
		jMap[k] = a.jMap[k]
	}
	b = &Args{
		iMap: iMap,
		jMap: jMap,
	}
	return
}
