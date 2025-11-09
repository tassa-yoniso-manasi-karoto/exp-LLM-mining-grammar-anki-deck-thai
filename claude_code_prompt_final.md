# Claude Code Project: JLab-Style Thai Anki Deck Creation - Phase 2

## Project Overview

Create a structured Thai Anki deck following the JLab methodology. The deck uses authentic Thai sentences from media (with audio, transcriptions, and translations) to teach intermediate grammar systematically.

**Key Resources:**
- `thai_deck_curriculum.md` - Complete curriculum with 1,500 cards organized into modules
- `tsv-finder-cli` - Command-line tool for searching Thai sentences (Phase 1 ✅ COMPLETE)
- `thai-sentence-selector` - Sub-agent for selecting best sentences for each grammar point
- TSV files - Sentences from Thai media with translations and audio

## PHASE 1: TSV Finder CLI Tool ✅ COMPLETED

A non-interactive CLI tool (`tsv-finder-cli`) has been developed with the following capabilities:

### Basic Usage
```bash
# Search and view results (default: Thai text with index)
./tsv-finder-cli --query "ให้" --native-only

# Output:
0	โอ๊ย แม่ไม่ต้องทำแบบนี้ละ...
1	แล้วตัวมึงเองก็จะเข้าใกล้เอเชียไม่ได้ด้วย...
2	หรือเสี่ยโต้งมันได้เงินไปแล้วครับ...
```

### Two-Stage Export Workflow
```bash
# Stage 1: Query to see results
./tsv-finder-cli --query "ให้" --native-only

# Stage 2: Export selected sentences
./tsv-finder-cli --query "ให้" --native-only --export-indices "0,4,7,12,15" --output-dir "[choosen working dir]"
# Output: [choosen working dir]/tsv_export_20251106_174321.tsv
```

### Export File Format
The `--export-indices` flag creates a TSV file with exactly 6 columns:
```tsv
[Audio]	[Timestamp]	[Episode]	[Image]	[Thai]	[English]
[sound:file.ogg]	00:44:22,285	Oh My Ghost_S01E07	"""<img src=""file.avif"">"""	แม่ซื้อของขวัญให้ฉัน	Mom bought a gift for me
```

**Columns:**
1. Audio file reference (with [sound:] tags)
2. Timestamp (HH:MM:SS,mmm format)
3. Episode/source name
4. Image HTML tag
5. Thai sentence
6. English translation

**Note:** Context sentences (previous/next) are excluded from exports.

## PHASE 2: Anki Deck Creation (CURRENT FOCUS)

### Deck File Management

**Deck TSV Location:** `./deck/thai_anki_deck.tsv`

**File Structure:**
- Cards are appended incrementally as each module is processed
- Both main agent and sub-agent have read access
- Only main agent writes to the file (append-only)
- **NO SCRIPTS OR LIBRARIES** - Direct TSV writing only
- **SIZE WARNING:** Deck will grow large (>25,000 tokens). Sub-agent must use `awk` to extract specific columns

**Initial Setup:**
```bash
# Create deck directory and file with tab-separated headers
mkdir -p ./deck
echo -e "Front\tThai\tPhonetic\tTranslation\tGrammar\tVocabulary\tAudio\tModule" > ./deck/thai_anki_deck.tsv
```

**During Processing:**
- Main agent directly appends TSV lines after each card generation
- Sub-agent reads current state to understand taught concepts
- File persists through the entire deck creation process
- NO Python scripts, NO external libraries - just direct file operations

### Sub-Agent Integration

**Using the `thai-sentence-selector` agent:**

The sub-agent expects these inputs:
- **Grammar point** to teach (e.g., "benefactive ให้")
- **Module context** (e.g., "Module 1.2, Cards 151-170")
- **Deck file path** to read existing cards
- **Search instructions** (what patterns to look for)

**How to Call the Sub-Agent:**
```
[Use the Task tool to invoke thai-sentence-selector with:]

Grammar Point: "Benefactive ให้ - indicates action done for someone else"
Module: "1.2, Cards 151-170 (20 cards needed)" ← IMPORTANT: Always specify exact number
Deck File: "./deck/thai_anki_deck.tsv"
Working Dir: Create your own in /tmp/ to avoid file conflicts
Search Target: "Find sentences with benefactive ให้ usage, avoiding causative or main verb uses"
```

**Note:** The sub-agent can select 10-40 cards per call. For larger sections, make multiple calls with different pattern variations.

**Sub-Agent's Process:**
1. Create its own isolated working directory
2. Extract vocabulary from deck using `awk` (avoids token limits):
   - `awk -F'\t' 'NR>1 {print $2}' /home/voiduser/go/src/tsv-finder/deck/thai_anki_deck.tsv` for Thai sentences
3. Construct and execute appropriate TSV finder queries
4. Analyze results for pedagogical value
5. Select the requested number of sentences based on i+1 principle
6. Export using `--export-indices --output-dir "[choosen working dir]"`
7. If multiple exports, concatenate carefully within working directory only
8. Return final TSV export filepath

**Sub-Agent Returns:**
- TSV file path (e.g., `/tmp/tsv_export_06174321.tsv`)

### Module Processing Workflow (Step-by-Step)

For each module in `thai_deck_curriculum.md`:

#### Step 1: Read Module Specifications
- Extract grammar point(s) from curriculum
- Note card count needed (e.g., "Cards 151-170" = 20 cards)
- Identify key patterns to teach

#### Step 2: Prepare Sub-Agent Call
- Formulate grammar point description
- Specify module and card range
- Describe what to search for (high-level)

#### Step 3: Call Sub-Agent
```
[Use Task tool to invoke thai-sentence-selector]
Grammar: "Benefactive ให้: V + Object + ให้ + Recipient"
Module: "1.2, Cards 151-170 (20 cards)"
Deck File: "./deck/thai_anki_deck.tsv"
Working Dir: Create your own in /tmp/ to avoid file conflicts
Search: "Find sentences with benefactive ให้, not causative or main verb"

Returns: TSV file path like /tmp/tsv_export_20251106_183045.tsv
```

#### Step 4: Process Exported Sentences
For each row in the returned TSV file:
1. Extract the 6 columns (Audio, Timestamp, Episode, Image, Thai, English)
2. Generate Paiboon romanization directly (no libraries - use your Thai knowledge AND read ./paiboon_examples.txt to get a refresher on the Paiboon transliteration)
3. Create grammar explanation following JLab principles:
   - **One New Thing:** Focus only on target grammar
   - **How, What, Why:** Structure, function, usage
   - **Anticipate Confusion:** Clarify potential issues
   - **Concise & Practical:** Brief and clear
4. Identify vocabulary needing notes (unfamiliar words only)
5. Format as tab-separated line

#### Step 5: Append to Deck
```bash
# Direct append to TSV file (no scripts!)
echo -e "[sound:file.ogg]\tไทย\tromanization\tEnglish\tgrammar\tvocab\t[sound:file.ogg]\t1.2" >> ./deck/thai_anki_deck.tsv
```

### Example Module Processing

**Module 1.1: The ก็ System (Cards 1-100)**

```
Step 1: Read module specifications from curriculum
- Basic ก็ patterns (Cards 1-30)
- ก็ + question words (Cards 31-70)
- Advanced ก็ expressions (Cards 71-100)

Step 2: Process each section sequentially

Section 1: Basic ก็ patterns
- Call sub-agent with:
  Grammar: "Basic ก็ patterns - topic continuation and result marking"
  Module: "1.1, Cards 1-30"
  Deck File: "./deck/thai_anki_deck.tsv"
  Working Dir: Create your own unique working directory following agent instructions
  Search: "Find ก็ in simple statements and conditions"

- Receive TSV export path
- Read each sentence from export
- Generate Paiboon romanization manually
- Write grammar explanation
- Append directly to deck TSV

Section 2: ก็ + question words
- Call sub-agent with:
  Grammar: "ก็ with question words - 'whatever' patterns"
  Module: "1.1, Cards 31-70"
  Deck File: "./deck/thai_anki_deck.tsv"
  Working Dir: Create your own unique working directory following agent instructions
  Search: "Find 'whatever' patterns like อะไรก็ได้, ที่ไหนก็ได้"

[Continue same process for Section 3]
```

### Card Generation Template

**Input from TSV:**
```tsv
[sound:show.ogg]	00:03:24	Show_Ep1	<img src="show.jpg">	แม่ซื้อของขวัญให้ฉัน	Mom bought a gift for me
```

**Generated Anki Card (TSV):**
```tsv
[sound:show.ogg]	แม่ซื้อของขวัญให้ฉัน	mɛ̂ɛ sɯ́ɯ khɔ̌ɔng-khwǎn hâi chǎn	Mom bought a gift for me	Benefactive ให้: V + Object + ให้ + Recipient. Action done FOR someone.	แม่ (mɛ̂ɛ): mom; ของขวัญ (khɔ̌ɔng-khwǎn): gift	[sound:show.ogg]	1.2
```

**Important:** Romanization is generated manually using your Thai knowledge and some examples of ./paiboon_examples.txt - no external libraries!

## Phase 2 Initialization Checklist

Before starting module processing:

1. **Create deck TSV file in safe location:**
   ```bash
   mkdir -p ./deck
   echo -e "Front\tThai\tPhonetic\tTranslation\tGrammar\tVocabulary\tAudio\tModule" > ./deck/thai_anki_deck.tsv
   ```

2. **Verify tsv-finder-cli is accessible:**
   ```bash
   ./tsv-finder-cli --help
   ```

3. **Verify thai-sentence-selector sub-agent is available**

4. **Load curriculum:**
   - Read `thai_deck_curriculum.md`
   - Understand module structure
   - Process in order: 1.1 → 1.2 → 1.3 → ... → 3.5
   - NO PARSING SCRIPTS - just read and process directly

5. **Begin with Module 1.1 (ก็ system)**

## Progress Tracking

- After each module section, verify card count
- Current progress always visible in `./deck/thai_anki_deck.tsv`
- **BACKUP REGULARLY:** `cp ./deck/thai_anki_deck.tsv ./deck/backup_[NAME OF LAST MODULE COMPLETED].tsv`
- Can resume from last module if interrupted
- Target: 1,500 total cards across all modules

## Critical Reminders

1. **Native content only** - Always use `--native-only` flag
2. **i+1 Principle** - Each card introduces ONE new concept
3. **Audio-first** - Front of card MUST be audio only
4. **No manual "previously taught" tracking** - Sub-agent reads deck file
5. **Append-only** - Never overwrite existing cards in deck file
6. **NO SCRIPTS OR LIBRARIES** - Direct TSV operations only
7. **Manual transliteration** - Use your Thai knowledge to generate Paiboon romanization directly and read ./paiboon_examples.txt to get a refresher on the Paiboon transliteration

## Next Steps

1. Initialize deck file with headers
2. Process Module 1.1 (ก็ system) - 100 cards
3. Verify card quality and format
4. Continue with remaining modules systematically
5. Final validation of complete 1,500-card deck