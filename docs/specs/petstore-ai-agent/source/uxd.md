# Petstore AI Agent UXD

## User Experience Design Document

**Project Name:** Petstore AI Recommendation Agent
**Author:** Design Team
**Date:** 2026-07-06
**Version:** 1.0
**Status:** Draft

## Design Principles

1. **Conversational, not transactional**: Feel like chatting with a knowledgeable friend
2. **Progressive disclosure**: Don't overwhelm; reveal complexity gradually
3. **Always helpful**: Even when declining, offer alternatives
4. **Transparent AI**: Be clear about capabilities and limitations

## User Personas

### Persona 1: Sarah, First-Time Pet Owner

**Demographics:** 28, urban apartment, works hybrid
**Goals:** Find a pet that fits her lifestyle
**Frustrations:** Overwhelmed by options, afraid of making wrong choice
**Behavior:** Asks many questions, needs reassurance

### Persona 2: Marcus, Busy Parent

**Demographics:** 42, suburban house, two kids (8, 12)
**Goals:** Find a kid-friendly, low-maintenance pet
**Frustrations:** Limited time for research, needs quick answers
**Behavior:** Direct questions, values efficiency

### Persona 3: Elena, Breed Enthusiast

**Demographics:** 35, experienced dog owner, knows what she wants
**Goals:** Check availability of specific breeds
**Frustrations:** Out-of-stock items, unclear availability
**Behavior:** Specific queries, expects quick inventory info

## User Journeys

### Journey 1: Discovery (Sarah)

```
Trigger: Sarah lands on homepage, feels overwhelmed

1. Notices chat bubble with "Need help finding your perfect pet?"
2. Clicks to open chat
3. Types: "I'm looking for my first pet"
4. Agent: "Welcome! I'd love to help. To find your perfect match,
          tell me about your living situation - apartment or house?"
5. Sarah: "Apartment, about 600 sq ft"
6. Agent: "Got it! Do you have a yard or balcony?"
7. Sarah: "Small balcony"
8. Agent: "Perfect. How many hours are you typically away from home?"
9. Sarah: "About 8-9 hours for work"
10. Agent: "Last question - any preference between dogs, cats, or
           open to smaller pets?"
11. Sarah: "Open to anything really"
12. Agent: "Based on your apartment lifestyle and work schedule,
           here are my top picks:

           🐱 **Persian Cat - $350**
           Perfect for apartments, independent, low exercise needs
           [View Whiskers →]

           🐰 **Holland Lop Rabbit - $120**
           Quiet, cage-friendly, great for small spaces
           [View Cinnamon →]

           🐠 **Betta Fish Tank Setup - $85**
           Zero noise, minimal space, beautiful display
           [View Options →]

           Would you like to know more about any of these?"
13. Sarah clicks "View Whiskers" and sees the cat detail page
14. Returns to chat: "What about grooming for the Persian?"
15. Agent provides grooming info with context maintained
```

### Journey 2: Availability Check (Elena)

```
Trigger: Elena wants a Golden Retriever puppy

1. Opens chat (returning user)
2. Types: "Do you have any Golden Retriever puppies?"
3. Agent: "Let me check our current Golden Retriever inventory...

          We have 2 Golden Retrievers available:

          🐕 **Max - Male, 4 months - $1,200**
          Vaccinated, started crate training
          [View Max →]

          🐕 **Daisy - Female, 6 months - $1,100**
          Spayed, knows basic commands
          [View Daisy →]

          Would you like details on either?"
4. Elena: "Any younger puppies coming soon?"
5. Agent: "We're expecting a litter in about 3 weeks.
          Would you like me to notify you when they arrive?"
6. Elena: "Yes please"
7. Agent: "Great! I'll email you at elena@email.com when
          Golden Retriever puppies become available.
          Anything else I can help with?"
```

### Journey 3: Compatibility Check (Marcus)

```
Trigger: Marcus's kids want a specific dog they saw

1. Marcus is on dog detail page for "Rocky" (German Shepherd)
2. Clicks "Ask AI about this pet"
3. Chat opens with context: "Asking about Rocky (German Shepherd)"
4. Marcus: "Is this dog good with kids ages 8 and 12?"
5. Agent: "Great question about Rocky for your family!

          ✅ **Good news:**
          - German Shepherds are known for loyalty to family
          - Rocky has been socialized with children
          - Great for active kids who like to play

          ⚠️ **Consider:**
          - They need 2+ hours of exercise daily
          - Can be protective; needs proper introduction
          - Requires consistent training

          Rocky specifically: He's 2 years old and has lived
          with a family before. He's good with kids but can
          be energetic during play.

          Would your family be able to commit to the exercise needs?"
6. Marcus: "We have a big backyard and the kids are active"
7. Agent: "That sounds like a great match! Rocky would love
          having space to run with the kids.

          One more thing - do you have other pets at home?"
8. Marcus: "A cat"
9. Agent: "Rocky's profile shows he's been tested with cats
          and does well with slow introductions.

          I'd suggest a meet-and-greet visit so your cat
          and Rocky can meet in a controlled environment.

          Would you like to schedule that?"
```

## Interface Design

### Chat Widget

```
┌─────────────────────────────────────┐
│  🐾 Pet Assistant              ─ × │
├─────────────────────────────────────┤
│                                     │
│  Agent message bubble (left)        │
│  ┌─────────────────────────┐        │
│  │ Hello! I'm here to help │        │
│  │ you find your perfect   │        │
│  │ pet companion.          │        │
│  └─────────────────────────┘        │
│                                     │
│         User message (right)        │
│        ┌─────────────────────┐      │
│        │ I need a small dog  │      │
│        └─────────────────────┘      │
│                                     │
│  Pet recommendation card            │
│  ┌─────────────────────────┐        │
│  │ 🐕 Buddy - Pomeranian   │        │
│  │ $450 • Male • 1 year    │        │
│  │ Perfect for apartments  │        │
│  │ [View Details →]        │        │
│  └─────────────────────────┘        │
│                                     │
├─────────────────────────────────────┤
│ ┌─────────────────────────────┐ 📤 │
│ │ Type your message...        │    │
│ └─────────────────────────────┘    │
└─────────────────────────────────────┘
```

### Visual Specifications

**Chat Bubble**
- Position: Fixed, bottom-right, 24px margin
- Size: 60px × 60px circle
- Icon: Paw print
- Animation: Gentle pulse every 30s when unopened

**Chat Window**
- Size: 380px wide × 520px tall (desktop)
- Mobile: Full-width, 70% height
- Background: White (#FFFFFF)
- Border radius: 12px
- Shadow: 0 4px 24px rgba(0,0,0,0.15)

**Typography**
- Agent messages: 14px, system font, #333333
- User messages: 14px, system font, white on #4A90D9
- Timestamps: 11px, #888888

**Colors**
- Primary: #4A90D9 (trustworthy blue)
- Success: #34C759 (available)
- Warning: #FF9500 (caution)
- Error: #FF3B30 (problem)
- Background: #F5F5F5 (chat area)

### Interaction States

**Loading State**
- Three-dot typing indicator
- Duration: Maximum 5 seconds visible
- Fallback: "Still thinking..." message

**Error State**
- Red banner at top of chat
- Retry button
- Never lose conversation history

**Empty State**
- Welcoming illustration
- "Hi! Ask me anything about pets"
- Suggested prompts as chips

### Suggested Prompts

Display as horizontal scrollable chips when chat is empty:

```
[Find my perfect pet] [Check availability] [Pet care info] [Compare pets]
```

## Conversation Design

### Agent Personality

- **Name:** "Pet Assistant" (no fake human name)
- **Tone:** Friendly, knowledgeable, enthusiastic about pets
- **Vocabulary:** Simple, avoid jargon, occasional emoji for warmth
- **Limitations:** Honest about what it can/can't do

### Response Patterns

**Greeting (first message)**
```
Hi! I'm your Pet Assistant. I can help you:
• Find pets that match your lifestyle
• Check availability of specific breeds
• Answer questions about pet care

What brings you to Petstore today?
```

**Clarification needed**
```
I want to make sure I understand - when you say [X],
do you mean [option A] or [option B]?
```

**Out of scope**
```
That's a great question, but it's outside my expertise.
For [veterinary advice/returns/etc.], please contact
our support team at support@petstore.example.com.

Is there anything about pet selection I can help with?
```

**No results**
```
I couldn't find any [X] in our current inventory.

Here's what I can do:
1. Notify you when [X] becomes available
2. Show you similar pets you might like

What would you prefer?
```

### Conversation Boundaries

**Will do:**
- Pet recommendations based on lifestyle
- Inventory and availability checks
- Care requirement explanations
- Pet temperament information
- Price comparisons

**Won't do:**
- Medical/veterinary advice
- Price negotiation
- Order placement (link to cart instead)
- Personal information requests beyond email for notifications
- Guarantees about pet behavior

## Accessibility

### WCAG 2.1 AA Compliance

- Color contrast: 4.5:1 minimum for text
- Focus indicators: Visible keyboard focus
- Screen reader: ARIA labels for all interactive elements
- Font size: 16px minimum, respects system preferences

### Keyboard Navigation

| Key | Action |
|-----|--------|
| Tab | Move between interactive elements |
| Enter | Send message, activate buttons |
| Escape | Close chat window |
| Up/Down | Navigate conversation history |

### Screen Reader Announcements

- "New message from Pet Assistant: [content]"
- "Pet recommendation: [name], [price], [description]"
- "Chat opened" / "Chat minimized"

## Error Handling

### Error Types and Responses

| Error | User Message | Action |
|-------|-------------|--------|
| Network timeout | "Having trouble connecting. Please try again." | Retry button |
| API error | "Something went wrong. Your message was saved." | Auto-retry |
| Rate limit | "Slow down! Let me catch up." | Disable input briefly |
| Invalid input | "I didn't understand that. Could you rephrase?" | Clear, try again |

### Graceful Degradation

If AI service unavailable:
1. Show static FAQ links
2. Offer direct contact options
3. Log for monitoring

## Analytics Events

| Event | Trigger | Data |
|-------|---------|------|
| chat_opened | User opens chat | page_url, user_segment |
| message_sent | User sends message | message_length, turn_number |
| recommendation_shown | Agent shows pets | pet_ids, query_type |
| recommendation_clicked | User clicks pet link | pet_id, position |
| notification_signup | User requests alerts | pet_type, email_hashed |
| chat_closed | User closes chat | conversation_length, last_intent |

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-07-06 | Design Team | Initial draft |
