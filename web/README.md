# Algoholic Web - Next.js Frontend

This is the Next.js 14 frontend for Algoholic, a LeetCode-style DSA practice platform with AI-powered learning paths.

## Tech Stack

- **Next.js 14** with App Router
- **TypeScript** for type safety
- **Tailwind CSS** for styling
- **React Query** for server state management
- **Zustand** for client state (auth)
- **Axios** for API requests
- **Lucide React** for icons

## Getting Started

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create `.env.local` file:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:4000/api
   ```

3. Run the development server:
   ```bash
   npm run dev
   ```

4. Open [http://localhost:3000](http://localhost:3000)

## Project Structure

```
src/
├── app/                    # Next.js App Router pages
│   ├── login/             # Login page
│   ├── register/          # Registration page
│   └── dashboard/         # Protected dashboard pages
│       ├── page.tsx       # Dashboard home
│       ├── problems/      # Problems browser
│       ├── practice/      # Practice session
│       ├── training-plans/# Training plans management
│       ├── lists/         # User problem lists
│       └── profile/       # User profile & activity
├── components/
│   ├── Providers.tsx      # React Query + Toast providers
│   └── DashboardLayout.tsx# Main layout wrapper
├── hooks/
│   └── useAuth.ts         # Auth hook
├── lib/
│   └── api.ts             # API client & types
└── stores/
    └── authStore.ts       # Zustand auth store
```

## Features

### Dashboard
- User stats overview (attempts, accuracy, streak)
- Strong/weak topic analysis
- Personalized recommendations
- Spaced repetition review queue

### Practice
- Random question practice
- Hint system with tracking
- Timed sessions
- Immediate feedback with explanations

### Problems
- Filter by difficulty, pattern, topic
- Semantic search
- View problem details and similar problems

### Training Plans
- Create custom study plans
- Topic/pattern focus options
- Adaptive difficulty setting
- Daily question goals

### Lists
- Create custom problem lists
- Track progress per list
- Public/private visibility

### Profile
- GitHub-style activity chart
- Practice history
- Streak tracking

## API Endpoints Used

All endpoints are aligned with the backend at `NEXT_PUBLIC_API_URL`:

- **Auth**: `/auth/register`, `/auth/login`, `/auth/me`
- **Problems**: `/problems`, `/problems/:id`, `/search/problems`
- **Questions**: `/questions/random`, `/questions/:id/answer`, `/questions/:id/hint`
- **Users**: `/users/me/stats`, `/users/me/recommendations`, `/users/me/review-queue`
- **Training Plans**: `/training-plans`, `/training-plans/:id/next`
- **Topics**: `/topics`, `/topics/:id/prerequisites`
- **Activity**: `/activity/chart`, `/activity/stats`

## Scripts

- `npm run dev` - Development server
- `npm run build` - Production build
- `npm run start` - Production server
- `npm run lint` - ESLint check
