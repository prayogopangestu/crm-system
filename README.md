# CRM Enterprise System

A modern, web-based Customer Relationship Management (CRM) system designed to help sales teams manage contacts, monitor sales pipelines, track daily tasks, and analyze business performance with a highly premium, fast, and interactive user interface.

---

## Key Features

This CRM application comes equipped with multiple ready-to-use modules styled using a modern Indigo & Slate theme:

### 1. Dashboard Overview
* **Bento Grid Layout**: A modern, responsive Bento-style layout for key information.
* **Core KPIs**: Metric cards showcasing Total Deals, Active Contacts, Completed Tasks, and Win Rate.
* **Interactive Charts (Recharts)**: An interactive area chart with premium color gradients for tracking revenue or activity trends.
* **Recent Activity Feed**: Real-time activity logs to monitor data modifications.

### 2. Contacts & Companies
* **Client Management**: A clean contact table with adaptive status pills (*Negotiation, Won, Lead, Proposal, Lost, Qualification*).
* **Contact Profiles**: Sleek initial-based avatars, email, job role, and last contacted date/time.

### 3. Sales Pipeline (Kanban Board)
* **Deal Stage Visualization**: An interactive Kanban board that divides deals into stages: *New Lead, Qualification, Negotiation, and Deal Won*.
* **Priority & Value Status**: Displays deal value (Rupiah/Currency) and priority tags (*High, Medium, Low*) along with the assignee.

### 4. Task & Activity Schedule (Task Management)
* **Interactive Calendar**: A dynamic monthly calendar with month navigation and filter functionality based on the selected date.
* **Task Categorization**: Tasks are automatically grouped into:
  * **Overdue**: Accompanied by a red badge count for urgent tasks.
  * **Today**: Today's task checklist showing schedule times.
  * **Upcoming**: Future task schedules.
* **3D Unrolling Paper Animation**: Clicking a task row expands it downwards using **Framer Motion** with a 3D unrolling paper effect (`rotateX` 3D, `transformOrigin: "top"`, parchment background, and a page-fold curl in the top-right corner).
* **Aligned & Synced Input Form**: A task creation form that supports Title, Related Client/Company, Activity Type (Call, Meeting, Proposal, Other), Time, Priority, Assignee, Detailed Notes, and a toggle to directly save the task as a completed logged activity.

### 5. Reports & Analytics
* **Donut Chart (Recharts)**: Visualization of "Deal Lost Reasons" to help make strategic business decisions.
* **Goal Tracking**: A comparative table of Monthly Goal vs. Actual Achievement with progress percentages.

### 6. System Settings
* Tabs panel to configure user profile details, app preferences, team members, customization of stages, and Telegram Webhook integration.

---

## Tech Stack

The application is built using the following modern stack:
* **Framework**: [Next.js 15+](https://nextjs.org/) (App Router & React 19)
* **Styling**: Vanilla CSS & [Tailwind CSS v4](https://tailwindcss.com/)
* **Animations**: [Framer Motion](https://www.framer.com/motion/)
* **Data Visualization**: [Recharts](https://recharts.org/)
* **UI Components**: Custom Radix-UI primitives & Lucide Icons
* **Notifications**: [Sonner Toast](https://sonner.emilkowal.ski/)

---

## Getting Started

### Prerequisites
Make sure you have Node.js version 18.x or later installed on your machine.

### 1. Clone the Repository
```bash
git clone <your-repository-url>
cd crm-system
```

### 2. Install Dependencies
Install all the required npm packages:
```bash
npm install
```

### 3. Run the Development Server
Start the local development server:
```bash
npm run dev
```
Open [http://localhost:3000](http://localhost:3000) in your browser to view the application.

### 4. Build for Production
To build the application for production deployment:
```bash
npm run build
npm run start
```

---

## Directory Structure

```text
crm-system/
├── src/
│   ├── app/                 # Next.js App Router (Pages & Layouts)
│   │   ├── tugas/           # Tasks & Calendar Module
│   │   ├── kontak/          # Client Contacts Module
│   │   ├── laporan/         # Analytics & Donut Chart Module
│   │   ├── pengaturan/      # Settings Module
│   │   ├── page.tsx         # Dashboard Overview
│   │   └── layout.tsx       # Global Layout
│   ├── components/          # Reusable UI & Layout Components
│   │   ├── common/          # Global Sidebar & Header
│   │   └── ui/              # Button, Card, Checkbox, Input, Select
│   ├── lib/                 # Helper utilities and Mock Data
│   └── styles/              # Global CSS configuration
```

---

## License
This project is licensed under the MIT License.
