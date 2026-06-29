import { create } from "zustand"
import { apiRequest, MutationResponse } from "@/lib/api"
import { Task, TaskPriority, TaskType } from "@/types/crm"

export interface TaskFormData {
  title: string
  relatedTo: string
  type: TaskType
  time: string
  priority: TaskPriority
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
  isLoading: boolean
  error: string | null
  
  // Task List Actions
  loadTasks: () => Promise<void>
  toggleTask: (taskId: string) => Promise<void>
  toggleExpand: (taskId: string) => void
  setCurrentMonthDate: (date: Date | ((prev: Date) => Date)) => void
  setSelectedDate: (date: Date) => void
  addTask: (task: {
    title: string
    company: string
    time: string
    date: string
    type: TaskType
    priority: TaskPriority
    assignee: string
    notes: string
    completed: boolean
  }) => Promise<void>
  
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
  assignee: "",
  completedDirectly: false,
  notes: ""
}

export const useTaskStore = create<TaskState>((set, get) => ({
  // Initial values
  tasks: [],
  currentMonthDate: new Date(),
  selectedDate: new Date(),
  expandedTaskId: null,
  step: 1,
  formData: initialFormData,
  isLoading: false,
  error: null,
  
  // Actions
  loadTasks: async () => {
    set({ isLoading: true, error: null })
    try {
      const tasks = await apiRequest<Task[]>("/api/tasks")
      set({ tasks, isLoading: false })
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal memuat tugas",
      })
    }
  },
  toggleTask: async (taskId) => {
    const previous = get().tasks
    const target = previous.find((task) => task.id === taskId)
    if (!target) return
    const completed = !target.completed
    set({
      tasks: previous.map((task) => (task.id === taskId ? { ...task, completed } : task)),
      error: null,
    })
    try {
      await apiRequest<{ success: boolean; completed: boolean }>(`/api/tasks/${taskId}/toggle`, {
        method: "PATCH",
        body: { completed },
      })
    } catch (error) {
      set({
        tasks: previous,
        error: error instanceof Error ? error.message : "Gagal mengubah status tugas",
      })
      throw error
    }
  },
  toggleExpand: (taskId) => set((state) => ({
    expandedTaskId: state.expandedTaskId === taskId ? null : taskId
  })),
  setCurrentMonthDate: (date) => set((state) => ({
    currentMonthDate: typeof date === "function" ? date(state.currentMonthDate) : date
  })),
  setSelectedDate: (selectedDate) => set({ selectedDate }),
  addTask: async (task) => {
    set({ isLoading: true, error: null })
    try {
      const response = await apiRequest<MutationResponse<Task>>("/api/tasks", {
        method: "POST",
        body: task,
      })
      set((state) => ({
        tasks: [response.data, ...state.tasks],
        isLoading: false,
      }))
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal menyimpan tugas",
      })
      throw error
    }
  },
  setStep: (step) => set({ step }),
  updateFormData: (data) => set((state) => ({
    formData: { ...state.formData, ...data }
  })),
  resetForm: () => set({ step: 1, formData: initialFormData })
}))
