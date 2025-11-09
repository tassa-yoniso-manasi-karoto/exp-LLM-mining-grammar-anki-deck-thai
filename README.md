# Experimental LLM Mining for a Thai A2+ Grammar Deck through Claude Code

Discussion: https://forums.ankiweb.net/t/can-claude-code-opus-make-high-quality-language-learning-decks-i-tried/67384

Deck available on Ankiweb: https://ankiweb.net/shared/info/2013132445

## Project Overview

This repository documents what is *AFAIK* the **first published attempt** at using Large Language Models to **autonomously create** a complete language-learning grammar deck. Through a **2-layer agent hierarchy** (lead Claude Code + sub-agents, both running Opus), Claude Opus 4.1 successfully generated **~1500 Thai grammar cards** (A2+ level) by mining authentic sentences from a **124,000-sentence corpus** extracted from Thai TV shows. The project demonstrates that while LLMs currently **cannot maintain perfect i+1** (incremental difficulty) progression without sophisticated scaffolding, they **can produce pedagogically useful grammar decks** with approximate difficulty ordering.

This opens possibilities for automated creation of near-human quality language learning notes for underserved languages.

# TLDR

### **Is it possible for a state-of-the-art LLM to create language learning decks following the i+1 sentences design?**

⟶ With **simple scaffolding** (1 CLI just to query for sentences, no programmatic pre-scoring of sentences) and a **2-layer agent hierarchy** (lead Claude Code + sub-agents), the answer is a **categoric NO**.

If more sophisticated scaffolding and a nested agent hierarchy (i.e. sub-agents that can create sub-sub-agents to lessen their workload) were used, the answer is ***perhaps***… but this is not supported in Claude Code for now.

### Are LLMs capable of writing a good language learning deck?

⟶ The answer is **yes**! They absolutely can do a good enough job for the resulting deck to be useful. In my case the grammar progression and the exhaustive showcasing of grammar points through real world examples looks great so I will be studying this Thai grammar deck.


## Methodology

I first created a **detailed Thai grammar (A2+) curriculum** representing about **1500 notes** to make. I passed this curriculum around between Claude Opus 4.1 and Gemini Pro 2.5, and I had them make refinements based on grammar books, so it eventually came to a pretty good and precise curriculum: `thai_deck_curriculum.md`

This curriculum served as the reference for actually writing the deck.
The setup was basically: **Claude Code acted as a main / lead agent** orchestrating the deck creation by using **sub-agents** to call the `tsv-finder-cli`, search through the sentence database results returned by the CLI, select the good ones that are relevant for the specific grammar point involved. That way, the task of going through the search or query results would be **delegated to the sub-agent** and the main lead agent, i.e. the Claude Code session would **not have its context window flooded** with the results.

Once the prompts were refined enough, which took a while, the process of making the deck was really relatively easy: just let the lead agent run with **`--dangerously-skip-permissions`** at full speed **without requesting any user confirmation**.
The main problem is that **user supervision is still required** to some extent because the prompt is quite elaborate and sometimes Claude Code will forget some important aspects of work to be done: especially when **restarting with a new lead agent session**.

In total, it took maybe **5 sliding windows of 5 hours** (that is the way Anthropic plans are rate limited). It consumed about **3.9M Opus tokens** and I was surprised by how difficult it was to actually reach the five-hour limits. But this is probably due to the fact that since it's not a coding scenario, most of the tokens involved are ***input tokens***.
And just in the same way that the API charges much less for input tokens, probably the Claude Code Max plan attributes much less usage consumption to input tokens. 

# Token Usage Summary
See `thai_deck_curriculum.md` for curriculum details.

## Token Usage Breakdown

<details>


### Lead Agent Sessions
- Session 1: **159k tokens** (ended after Module 1.3 Purpose & acquisition chains)
- Session 2: **158k tokens** (ended after Module 1.4 Duration & continuity)
- Session 3: **179k tokens** (ended after Module 2.5 Habitual patterns)
- Session 4: **183k tokens** (ended after Module 3.2 Colloquial & Regional)
- Session 5: **154k tokens** (ended after Module 3.4 การ structures)
- Session 6: **123k tokens** (ended after Module 3.5 Hesitation markers)

**Total Lead Agent Tokens: 956k**


## PHASE 1: CRITICAL GAPS (Cards 1-400)

### Module 1.1: Core ก็ (gɔ̂) System (Cards 1-50)
- Basic ก็ patterns (Cards 1-10) — **59.6k subagent tokens**
- Conditionals with ก็ (Cards 11-20) — **70.4k subagent tokens**
- Question word + ก็ได้ patterns (Cards 21-35) — **52.2k subagent tokens**
- Essential ก็ expressions (Cards 36-45) — **68.1k subagent tokens**
- Special uses (Cards 46-50) — *[not measured]*

### Module 1.2: Complete ให้ (hâi) System (Cards 51-150)
- Basic giving & receiving (Cards 101-120) — **87.5k subagent tokens**
- Causative ให้ (Cards 121-150) — **81.5k subagent tokens**
- Benefactive ให้ (Cards 151-170) — **55.7k subagent tokens**
- Purpose/extent ให้ (Cards 171-200) — **72.2k subagent tokens**
- Complex chains — **58.1k subagent tokens**

### Module 1.3: Serial Verb Mastery (Cards 151-300)
- Basic directional verbs (Cards 151-180) — **101.7k subagent tokens**
- Purpose & acquisition chains (Cards 181-210) — **80.8k subagent tokens**
- Helper & assistance chains (Cards 211-240) — **66.7k subagent tokens**
- Perception & discovery chains (Cards 241-260) — **91.8k subagent tokens**
- Result & completion chains (Cards 261-280) — **64.0k subagent tokens**
- Abstract & mental chains (Cards 281-300) — **65.7k subagent tokens**

### Module 1.4: Temporal Sequencing & Aspect (Cards 301-400)
- Aspect markers & combinations (Cards 301-340) — **81.7k subagent tokens**
- Temporal sequencing (Cards 341-370) — **85.0k subagent tokens**
- Duration & continuity (Cards 371-400) — **83.3k subagent tokens**

**PHASE 1 SUBAGENT TOTAL: ~1,346.2k tokens**

---

## PHASE 2: INTERMEDIATE MASTERY (Cards 401-1000)

### Module 2.1: Complete Particle System (Cards 401-500)
- ล่ะ (lâ) complete system (Cards 401-430) — **69.7k subagent tokens**
- Other essential particles (Cards 431-500) — **80.5k subagent tokens**

### Module 2.2: Advanced Conditionals (Cards 501-600)
- *[subsections not measured separately]*

### Module 2.3: Contrast & Concession (Cards 601-700)
- *[subsections not measured separately]*

### Module 2.4: Quantifiers & Degree (Cards 701-800)
- Scalar quantifiers (Cards 701-750) — **69.4k subagent tokens**
- Sufficiency & excess (Cards 751-800) — **104.0k subagent tokens**

### Module 2.5: Frequency & Habits (Cards 801-900)
- Frequency adverbs (Cards 801-850) — **71.1k subagent tokens**
- Habitual patterns (Cards 851-900) — **61.0k subagent tokens**

### Module 2.6: Voice & Perspective (Cards 901-1000)
- Passive constructions (Cards 901-950) — **66.0k subagent tokens**
- Impersonal & generic (Cards 951-1000) — **74.6k subagent tokens**

**PHASE 2 SUBAGENT PARTIAL TOTAL: ~596.3k tokens** *(some modules unmeasured)*

---

## PHASE 3: ADVANCED NATURAL THAI (Cards 1001-1500)

### Module 3.1: Discourse Management (Cards 1001-1100)
- Topic & focus markers (Cards 1001-1050) — **78.5k subagent tokens**
- Hedging & mitigation (Cards 1051-1100) — **84.6k subagent tokens**

### Module 3.2: Colloquial & Regional (Cards 1101-1200)
- Contractions & reductions (Cards 1101-1150) — **98.4k subagent tokens** *(combined measurement)*
- Regional variations & slang (Cards 1151-1200) — *[included above]*

### Module 3.3: Advanced Comparison & Appearance (Cards 1201-1300)
- Similarity & difference nuances (Cards 1201-1250) — **93.4k subagent tokens** *(combined measurement)*
- Appearance & epistemic modality (Cards 1251-1300) — *[included above]*

### Module 3.4: Complex Nominalizations (Cards 1301-1400)
- การ (gaan) structures (Cards 1301-1350) — **56.3k subagent tokens**
- ความ (khwaam) abstractions (Cards 1351-1400) — **85.5k subagent tokens**

### Module 3.5: Idiomatic Mastery (Cards 1401-1500)
- Fixed expressions (Cards 1401-1450) — **73.6k subagent tokens**
- Prosody & pragmatics (Cards 1451-1500):
  - Emphasis through lengthening — **77.9k subagent tokens**
  - Irony & sarcasm markers — **65.1k subagent tokens**
  - Turn-taking signals — **64.6k subagent tokens**
  - Backchanneling expressions — **106.4k subagent tokens**
  - Hesitation markers — **81.0k subagent tokens**

**PHASE 3 SUBAGENT TOTAL: ~965.3k tokens**

</details>


## Summary Statistics

### Subagent Tokens
- **Phase 1 (Cards 1-400):** ~1.35M tokens
- **Phase 2 (Cards 401-1000):** ~0.60M tokens (partial)
- **Phase 3 (Cards 1001-1500):** ~0.97M tokens
- **Subagent Total:** ~2.92M tokens

### Lead Agent Tokens
- **6 Sessions Total:** 956k tokens

### GRAND TOTAL: ~3.88M tokens

Calculated the cost for Claude Opus 4.1 with 3.88M tokens (assuming 90% input, 10% output):
* Input tokens: 3.88M × 0.90 = 3.492M tokens
* Output tokens: 3.88M × 0.10 = 0.388M tokens


According to the pricing documentation, Claude Opus 4.1 costs $15 per million input tokens and $75 per million output tokens, if no caching is used.


Cost breakdown:
* Input: 3.492M tokens × ($15 / 1M tokens) = $52.38
* Output: 0.388M tokens × ($75 / 1M tokens) = $29.10
### Estimated cost if API had been used: $81.48
