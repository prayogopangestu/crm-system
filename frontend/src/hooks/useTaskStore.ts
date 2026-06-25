import { create } from "zustand"
import { initialTasks, Task } from "@/lib/mockData"

export interface TaskFormData {
  title: string
  relatedTo: string
  type: "Meeting" | "Call" | "Proposal" | "Other"
  time: string
  priority: "Tinggi" | "Sedang" | "Rendah"
  assignee: string
  completedDirectly: boolean
  notes: string
}

interface TaskState {
  // Task List States
  tasks: Task[]
  currentMonthDate: Date
  selectedDate: Date
  expandedTaskId: string | null
  
  // Form States
  step: number
  formData: TaskFormData
  
  // Task List Actions
  toggleTask: (taskId: string) => void
  toggleExpand: (taskId: string) => void
  setCurrentMonthDate: (date: Date | ((prev: Date) => Date)) => void
  setSelectedDate: (date: Date) => void
  addTask: (task: Task) => void
  
  // Form Actions
  setStep: (step: number) => void
  updateFormData: (data: Partial<TaskFormData>) => void
  resetForm: () => void
}

const initialFormData: TaskFormData = {
  title: "",
  relatedTo: "",
  type: "Call",
  time: "12:00",
  priority: "Sedang",
  assignee: "Sarah Jenkins",
  completedDirectly: false,
  notes: ""
}

export const useTaskStore = create<TaskState>((set) => ({
  // Initial values
  tasks: initialTasks,
  currentMonthDate: new Date(),
  selectedDate: new Date(),
  expandedTaskId: null,
  step: 1,
  formData: initialFormData,
  
  // Actions
  toggleTask: (taskId) => set((state) => ({
    tasks: state.tasks.map((task) =>
      task.id === taskId ? { ...task, completed: !task.completed } : task
    )
  })),
  toggleExpand: (taskId) => set((state) => ({
    expandedTaskId: state.expandedTaskId === taskId ? null : taskId
  })),
  setCurrentMonthDate: (date) => set((state) => ({
    currentMonthDate: typeof date === "function" ? date(state.currentMonthDate) : date
  })),
  setSelectedDate: (selectedDate) => set({ selectedDate }),
  addTask: (task) => set((state) => ({
    tasks: [task, ...state.tasks]
  })),
  setStep: (step) => set({ step }),
  updateFormData: (data) => set((state) => ({
    formData: { ...state.formData, ...data }
  })),
  resetForm: () => set({ step: 1, formData: initialFormData })
}))
