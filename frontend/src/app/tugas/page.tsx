"use client"

import React, { useEffect } from "react"
import { Card } from "@/components/ui/Card"
import { Checkbox } from "@/components/ui/Checkbox"
import { Task } from "@/types/crm"
import { toast } from "sonner"
import { motion, AnimatePresence } from "framer-motion"
import { cn } from "@/lib/utils"
import { useTaskStore, TaskFormData } from "@/hooks/useTaskStore"
import { Step1Basic } from "@/components/tugas/Step1Basic"
import { Step2Config } from "@/components/tugas/Step2Config"
import { Step3Summary } from "@/components/tugas/Step3Summary"
import { useLanguage } from "@/context/LanguageContext"

export default function TugasPage() {
  const {
    tasks,
    currentMonthDate,
    selectedDate,
    expandedTaskId,
    step,
    error,
    loadTasks,
    toggleTask,
    toggleExpand,
    setCurrentMonthDate,
    setSelectedDate,
    addTask,
    resetForm
  } = useTaskStore()
  const { t, tList } = useLanguage()

  useEffect(() => {
    void loadTasks()
  }, [loadTasks])

  const months = tList("values.months")
  const days = tList("values.days")

  const handlePrevMonth = () => {
    setCurrentMonthDate(prev => {
      const nextDate = new Date(prev)
      nextDate.setMonth(prev.getMonth() - 1)
      return nextDate
    })
  }

  const handleNextMonth = () => {
    setCurrentMonthDate(prev => {
      const nextDate = new Date(prev)
      nextDate.setMonth(prev.getMonth() + 1)
      return nextDate
    })
  }

  const handleToday = () => {
    const today = new Date()
    setCurrentMonthDate(today)
    setSelectedDate(today)
    toast.success(t("tasks.todayRedirected"))
  }

  const renderTaskDetails = (task: Task) => {
    return (
      <AnimatePresence>
        {expandedTaskId === task.id && (
          <motion.div
            initial={{ opacity: 0, height: 0, rotateX: -35 }}
            animate={{ opacity: 1, height: "auto", rotateX: 0 }}
            exit={{ opacity: 0, height: 0, rotateX: -35 }}
            transition={{ 
              height: { type: "spring", stiffness: 150, damping: 22, mass: 0.9 },
              opacity: { duration: 0.25, ease: "easeInOut" },
              rotateX: { type: "spring", stiffness: 120, damping: 18, mass: 0.9 }
            }}
            className="overflow-hidden mt-3"
            style={{ 
              transformOrigin: "top center", 
              transformStyle: "preserve-3d", 
              perspective: "1200px" 
            }}
          >
            <div className="bg-amber-50/40 dark:bg-slate-900/50 border border-amber-100/50 dark:border-slate-800 rounded-lg p-3 text-[11px] shadow-inner relative text-left">
              {/* Curl page fold effect */}
              <div className="absolute top-0 right-0 w-8 h-8 bg-gradient-to-bl from-amber-200/20 to-transparent dark:from-slate-800/40 rounded-tr-lg pointer-events-none"></div>
              
              <div className="grid grid-cols-2 gap-3 mb-2 pb-2 border-b border-amber-100/30 dark:border-slate-800/50">
                <div>
                  <span className="text-slate-400 dark:text-slate-500 font-medium block">{t("tasks.priority")}</span>
                  <span className="font-semibold text-slate-700 dark:text-slate-300 flex items-center gap-1 mt-0.5">
                    <span className={`w-1.5 h-1.5 rounded-full ${
                      task.priority.includes("Tinggi")
                        ? "bg-red-500"
                        : task.priority.includes("Rendah")
                          ? "bg-slate-400 dark:bg-slate-500"
                          : "bg-amber-500"
                    }`}></span>
                    {t(`values.priority.${task.priority}`)}
                  </span>
                </div>
                <div>
                  <span className="text-slate-400 dark:text-slate-500 font-medium block">{t("tasks.assignee")}</span>
                  <span className="font-semibold text-slate-700 dark:text-slate-300 mt-0.5 block">{task.assignee || "-"}</span>
                </div>
              </div>

              <div>
                <span className="text-slate-400 dark:text-slate-500 font-medium block">{t("tasks.detailNotes")}</span>
                <p className="mt-1 text-slate-600 dark:text-slate-400 leading-relaxed font-normal normal-case">
                  {task.notes || t("tasks.noNotes")}
                </p>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    )
  }

  const getTaskDate = (task: Task): Date => {
    return new Date(`${task.date}T00:00:00`)
  }

  const isSameDay = (d1: Date, d2: Date) => 
    d1.getDate() === d2.getDate() && 
    d1.getMonth() === d2.getMonth() && 
    d1.getFullYear() === d2.getFullYear()

  // Save new logged activity as a completed task in state from Multi-step form
  const formatDate = (date: Date) => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, "0")
    const day = String(date.getDate()).padStart(2, "0")
    return `${year}-${month}-${day}`
  }

  const handleSaveMultiStepTask = async (data: TaskFormData) => {
    try {
      await addTask({
        title: data.title,
        company: data.relatedTo,
        time: data.time,
        date: formatDate(selectedDate),
        type: data.type,
        priority: data.priority,
        assignee: data.assignee,
        notes: data.notes,
        completed: data.completedDirectly,
      })

      toast.success(t("toast.taskAdded"), {
        description: t("toast.taskAddedDesc", { title: data.title, company: data.relatedTo, date: `${selectedDate.getDate()} ${months[selectedDate.getMonth()]}` }),
      })
      resetForm()
    } catch {
      toast.error(t("toast.taskSaveError"))
    }
  }

  // Filter tasks dynamically
  const today = new Date()
  const isFilteringByDate = !isSameDay(selectedDate, today)

  const getFilteredTasksForDate = (date: Date) => {
    return tasks.filter(t => isSameDay(getTaskDate(t), date))
  }

  // Standard lists for "Today" view
  const overdueTasks = tasks.filter(t => t.status === "overdue")
  const todayTasks = tasks.filter(t => t.status === "today")
  const upcomingTasks = tasks.filter(t => t.status === "upcoming")
  
  // Clicked date list
  const selectedDateTasks = isFilteringByDate ? getFilteredTasksForDate(selectedDate) : []

  // Generate calendar grid
  const getDaysInMonth = (date: Date) => {
    const year = date.getFullYear()
    const month = date.getMonth()
    
    const totalDays = new Date(year, month + 1, 0).getDate()
    const firstDayIndex = new Date(year, month, 1).getDay()
    const startOffset = firstDayIndex === 0 ? 6 : firstDayIndex - 1

    const days = []
    
    const prevMonthTotalDays = new Date(year, month, 0).getDate()
    for (let i = startOffset - 1; i >= 0; i--) {
      const dayNum = prevMonthTotalDays - i
      days.push({
        dayNum,
        isCurrentMonth: false,
        date: new Date(year, month - 1, dayNum)
      })
    }

    for (let i = 1; i <= totalDays; i++) {
      days.push({
        dayNum: i,
        isCurrentMonth: true,
        date: new Date(year, month, i)
      })
    }

    const totalCells = days.length > 35 ? 42 : 35
    const nextMonthDaysNeeded = totalCells - days.length
    for (let i = 1; i <= nextMonthDaysNeeded; i++) {
      days.push({
        dayNum: i,
        isCurrentMonth: false,
        date: new Date(year, month + 1, i)
      })
    }

    return days
  }

  const calendarDays = getDaysInMonth(currentMonthDate)

  return (
    <div className="p-gutter max-w-container-max-width mx-auto w-full flex flex-col gap-gutter lg:flex-row">
      {/* Left Column: Calendar & Quick Log */}
      <div className="flex-1 flex flex-col gap-gutter min-w-[60%] animate-fade-in-up">
        {error && (
          <div className="rounded-lg border border-error-container bg-error-container/40 px-4 py-3 text-xs font-semibold text-on-error-container">
            {error}
          </div>
        )}

        {/* Page Title & Navigation */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-end gap-4">
          <div>
            <h2 className="font-headline-md text-[24px] font-semibold text-on-surface">
              {t("tasks.title")}
            </h2>
            <p className="font-body-md text-sm text-on-surface-variant mt-1">
              {t("tasks.subtitle")}
            </p>
          </div>
          <div className="flex gap-2 self-end items-center">
            <span className="font-semibold text-sm mr-2 text-on-surface dark:text-slate-200">
              {months[currentMonthDate.getMonth()]} {currentMonthDate.getFullYear()}
            </span>
            <button 
              onClick={handlePrevMonth}
              className="p-2 rounded-lg bg-surface border border-outline-variant text-on-surface hover:bg-surface-variant transition-colors flex items-center justify-center cursor-pointer"
            >
              <span className="material-symbols-outlined text-[20px]">chevron_left</span>
            </button>
            <button
              onClick={handleToday}
              className="px-4 py-2 rounded-lg bg-surface border border-outline-variant text-on-surface hover:bg-surface-variant transition-colors font-label-md text-sm font-semibold cursor-pointer"
            >
              {t("tasks.today")}
            </button>
            <button 
              onClick={handleNextMonth}
              className="p-2 rounded-lg bg-surface border border-outline-variant text-on-surface hover:bg-surface-variant transition-colors flex items-center justify-center cursor-pointer"
            >
              <span className="material-symbols-outlined text-[20px]">chevron_right</span>
            </button>
          </div>
        </div>

        {/* Calendar Card (Bento style) */}
        <Card className="p-6 flex-1 flex flex-col">
          {/* Days Header */}
          <div className="grid grid-cols-7 gap-4 mb-4 border-b border-outline-variant pb-4">
            <div className="text-center font-label-sm text-xs font-semibold text-on-surface-variant uppercase">{days[0]}</div>
            <div className="text-center font-label-sm text-xs font-semibold text-on-surface-variant uppercase">{days[1]}</div>
            <div className="text-center font-label-sm text-xs font-semibold text-on-surface-variant uppercase">{days[2]}</div>
            <div className="text-center font-label-sm text-xs font-semibold text-on-surface-variant uppercase">{days[3]}</div>
            <div className="text-center font-label-sm text-xs font-semibold text-on-surface-variant uppercase">{days[4]}</div>
            <div className="text-center font-label-sm text-xs font-semibold text-on-surface-variant uppercase">{days[5]}</div>
            <div className="text-center font-label-sm text-xs font-semibold text-on-surface-variant uppercase">{days[6]}</div>
          </div>

          {/* Calendar Grid */}
          <div className="grid grid-cols-7 gap-2 flex-1 min-h-[300px]">
            {calendarDays.map((day, idx) => {
              const tasksForDay = getFilteredTasksForDate(day.date)
              const isTodayCell = isSameDay(day.date, today)
              const isSelectedCell = isSameDay(day.date, selectedDate)
              
              return (
                <div
                  key={idx}
                  onClick={() => {
                    setSelectedDate(day.date)
                    toast.info(t("tasks.viewingDate", { date: `${day.dayNum} ${months[day.date.getMonth()]} ${day.date.getFullYear()}` }))
                  }}
                  className={`p-2 rounded-lg border transition-all min-h-[80px] flex flex-col cursor-pointer select-none relative ${
                    isSelectedCell
                      ? "border-primary bg-primary-fixed/10"
                      : isTodayCell
                        ? "border-primary/50 bg-surface-container-low"
                        : day.isCurrentMonth
                          ? "border-outline-variant/40 bg-surface hover:border-primary/30"
                          : "border-transparent bg-slate-50/20 text-on-surface-variant opacity-40 hover:border-outline-variant"
                  }`}
                >
                  {isTodayCell && (
                    <span className="absolute top-1.5 right-1.5 w-1.5 h-1.5 bg-primary rounded-full"></span>
                  )}
                  <span className={`text-xs font-semibold mb-1 ${
                    isTodayCell ? "text-primary font-bold" : ""
                  }`}>
                    {day.dayNum}
                  </span>
                  
                  {/* Event Badges */}
                  <div className="flex-1 flex flex-col justify-end overflow-hidden">
                    {tasksForDay.slice(0, 2).map(task => (
                      <div
                        key={task.id}
                        className={`mt-0.5 text-[9px] px-1.5 py-0.5 rounded truncate font-medium ${
                          task.completed
                            ? "bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-500 line-through"
                            : task.status === "overdue"
                              ? "bg-red-500/10 text-red-700 dark:bg-red-950/20 dark:text-red-400"
                              : "bg-primary-fixed text-on-primary-fixed-variant"
                        }`}
                        title={task.title}
                      >
                        {task.title}
                      </div>
                    ))}
                    {tasksForDay.length > 2 && (
                      <div className="text-[8px] text-slate-400 dark:text-slate-500 font-bold mt-0.5 text-right">
                        {t("tasks.moreCount", { count: tasksForDay.length - 2 })}
                      </div>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        </Card>

        {/* Quick Log Form (Multi-step Form) */}
        <Card className="p-6">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-4 pb-3 border-b border-outline-variant/30">
            <h3 className="font-headline-sm text-base font-semibold text-on-surface flex items-center gap-2">
              <span className="material-symbols-outlined text-primary">add_task</span>
              {t("tasks.addTaskTitle", { date: `${selectedDate.getDate()} ${months[selectedDate.getMonth()]}` })}
            </h3>

            {/* Step Indicators */}
            <div className="flex items-center gap-1.5 text-[10px] font-semibold text-on-surface-variant">
              <span className={cn("px-2 py-0.5 rounded-full transition-all duration-200", step >= 1 ? "bg-primary text-on-primary" : "bg-surface-container-high text-on-surface-variant")}>
                {t("tasks.stepBasic")}
              </span>
              <span className="text-outline-variant">&rarr;</span>
              <span className={cn("px-2 py-0.5 rounded-full transition-all duration-200", step >= 2 ? "bg-primary text-on-primary" : "bg-surface-container-high text-on-surface-variant")}>
                {t("tasks.stepDetail")}
              </span>
              <span className="text-outline-variant">&rarr;</span>
              <span className={cn("px-2 py-0.5 rounded-full transition-all duration-200", step >= 3 ? "bg-primary text-on-primary" : "bg-surface-container-high text-on-surface-variant")}>
                {t("tasks.stepSummary")}
              </span>
            </div>
          </div>

          {/* Steps Orchestrator */}
          {step === 1 && <Step1Basic />}
          {step === 2 && <Step2Config />}
          {step === 3 && (
            <Step3Summary
              onSave={handleSaveMultiStepTask}
              onCancel={() => resetForm()}
            />
          )}
        </Card>
      </div>

      {/* Right Column: Task Lists */}
      <div className="w-full lg:w-[380px] flex flex-col gap-6">
        {isFilteringByDate ? (
          <Card className="p-4 flex flex-col relative overflow-hidden transition-all duration-300 hover:shadow-md hover:border-outline-variant dark:hover:border-slate-700 hover:-translate-y-[1px] animate-fade-in-up">
            <div className="flex items-center justify-between mb-4 border-b border-outline-variant/60 pb-3 relative z-10">
              <div>
                <h3 className="font-semibold text-sm text-slate-900 dark:text-slate-100">
                  {t("tasks.tasksForDate", { date: `${selectedDate.getDate()} ${months[selectedDate.getMonth()]}` })}
                </h3>
                <p className="text-[10px] text-slate-500 font-medium">{t("tasks.selectedDateSchedule")}</p>
              </div>
              <button
                onClick={() => setSelectedDate(today)}
                className="px-2.5 py-1 bg-primary text-white text-[10px] font-bold rounded-full cursor-pointer hover:bg-primary/95 transition-colors"
              >
                {t("tasks.today")}
              </button>
            </div>

            <div className="divide-y divide-outline-variant/30 relative z-10">
              {selectedDateTasks.length === 0 ? (
                <div className="py-12 text-center flex flex-col items-center justify-center">
                  <span className="material-symbols-outlined text-slate-400/50 text-[36px] mb-2">event_busy</span>
                  <p className="text-xs text-slate-500 font-semibold">{t("tasks.noTasksForDate")}</p>
                  <p className="text-[10px] text-slate-400 mt-1">{t("tasks.useFormLeft")}</p>
                </div>
              ) : (
                selectedDateTasks.map(task => (
                  <div 
                    key={task.id} 
                    onClick={() => toggleExpand(task.id)}
                    className={`py-3 flex flex-col gap-1 group relative transition-colors cursor-pointer select-none ${
                      task.completed ? "opacity-60" : "hover:bg-slate-50/40 dark:hover:bg-slate-800/10 px-2.5 -mx-2.5 rounded-lg"
                    }`}
                  >
                    <div className="flex items-start gap-3 w-full">
                      <div onClick={(e) => e.stopPropagation()}>
                        <Checkbox 
                          checked={task.completed} 
                          onChange={() => void toggleTask(task.id)} 
                          className="mt-1 cursor-pointer" 
                        />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-2">
                          <h4 className={`font-semibold text-sm text-slate-800 dark:text-slate-200 group-hover:text-primary transition-colors ${
                            task.completed ? "line-through text-slate-400 dark:text-slate-600" : ""
                          }`}>
                            {task.title}
                          </h4>
                          <span className={`shrink-0 px-2 py-0.5 rounded-full text-[9px] font-bold tracking-wider uppercase ${
                            task.completed 
                              ? "bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-600" 
                              : "bg-slate-100 text-slate-700 border border-slate-200/50 dark:bg-slate-800 dark:text-slate-300 dark:border-slate-750"
                          }`}>
                            {task.type}
                          </span>
                        </div>
                        <p className="text-xs text-slate-500 dark:text-slate-400 mt-1 truncate">{task.company}</p>
                        <div className="flex items-center justify-between mt-3 pt-2.5 border-t border-slate-100 dark:border-slate-800/60">
                          <div className="text-[10px] text-slate-400 font-semibold">
                            <span>{t("tasks.timeLabel", { time: task.time })}</span>
                          </div>
                          {task.completed ? (
                            <span className="text-[10px] text-emerald-600 font-bold">{t("tasks.completed")}</span>
                          ) : (
                            <span className="text-[10px] text-amber-600 font-bold">{t("tasks.notCompleted")}</span>
                          )}
                        </div>
                      </div>
                    </div>
                    {renderTaskDetails(task)}
                  </div>
                ))
              )}
            </div>
          </Card>
        ) : (
          <>
            {/* Overdue Tasks */}
            {overdueTasks.length > 0 && (
              <Card className="p-4 flex flex-col relative overflow-hidden transition-all duration-300 hover:shadow-md hover:border-outline-variant dark:hover:border-slate-700 hover:-translate-y-[1px] animate-fade-in-up">
                <div className="flex items-center justify-between mb-4 border-b border-outline-variant/60 pb-3 relative z-10">
                  <div>
                    <h3 className="font-semibold text-sm text-slate-900 dark:text-slate-100">
                      {t("tasks.overdueTitle")}
                    </h3>
                    <p className="text-[10px] text-slate-500 dark:text-slate-400 font-medium">{t("tasks.overdueSubtitle")}</p>
                  </div>
                  <span className="px-2.5 py-0.5 bg-red-500 text-white text-[10px] font-bold rounded-full shadow-xs shadow-red-500/10">
                    {t("tasks.urgent", { count: overdueTasks.length })}
                  </span>
                </div>

                <div className="divide-y divide-outline-variant/30 relative z-10">
                  {overdueTasks.map(task => (
                    <div 
                      key={task.id} 
                      onClick={() => toggleExpand(task.id)}
                      className={`py-3.5 first:pt-0 last:pb-0 flex flex-col gap-1 group relative transition-colors cursor-pointer select-none ${
                        task.completed ? "opacity-60" : "hover:bg-slate-50/40 dark:hover:bg-slate-800/10 px-2.5 -mx-2.5 rounded-lg"
                      }`}
                    >
                      <div className="flex items-start gap-3 w-full">
                        {/* Left accent indicator (rendered next to checkbox when row is active) */}
                        {!task.completed && (
                          <div className="w-[3px] self-stretch bg-red-500 rounded-full my-0.5"></div>
                        )}
                        
                        <div className="flex items-start gap-3 w-full">
                          <div onClick={(e) => e.stopPropagation()}>
                            <Checkbox
                              checked={task.completed}
                              onChange={() => {
                                void toggleTask(task.id)
                                if (!task.completed) {
                                  toast.success(t("toast.taskCompleted"), {
                                    description: t("toast.taskCompletedDesc", { title: task.title }),
                                  })
                                }
                              }}
                              className="mt-1 accent-red-600 cursor-pointer"
                            />
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="flex items-start justify-between gap-2">
                              <h4 className={`font-semibold text-sm text-slate-800 dark:text-slate-200 group-hover:text-primary transition-colors ${
                                task.completed ? "line-through text-slate-400 dark:text-slate-600" : ""
                              }`}>
                                {task.title}
                              </h4>
                              <span className={`shrink-0 px-2 py-0.5 rounded-full text-[9px] font-bold tracking-wider uppercase ${
                                task.completed 
                                  ? "bg-slate-100 text-slate-400 dark:bg-slate-800 dark:text-slate-600" 
                                  : "bg-slate-100 text-slate-700 border border-slate-200/50 dark:bg-slate-800 dark:text-slate-300 dark:border-slate-750"
                              }`}>
                                {task.type}
                              </span>
                            </div>
                            
                            <p className="text-xs text-slate-500 dark:text-slate-400 mt-1 truncate">
                              {task.company}
                            </p>
                            
                            <div className="flex items-center justify-between mt-3 pt-2.5 border-t border-slate-100 dark:border-slate-800/60">
                              <div className={`text-[10px] font-semibold ${
                                task.completed
                                  ? "text-slate-400 dark:text-slate-600"
                                  : "text-slate-600 dark:text-slate-350"
                              }`}>
                                <span>{t("tasks.dueTime", { time: task.time })}</span>
                              </div>

                              {!task.completed && (
                                <button
                                  onClick={(e) => {
                                    e.stopPropagation()
                                    void toggleTask(task.id)
                                    toast.success(t("toast.taskCompleted"), {
                                      description: t("toast.taskCompletedDesc", { title: task.title }),
                                    })
                                  }}
                                  className="text-[10px] text-slate-400 hover:text-primary dark:text-slate-500 dark:hover:text-primary font-bold flex items-center transition-colors cursor-pointer"
                                >
                                  {t("tasks.markDone")}
                                </button>
                              )}
                            </div>
                          </div>
                        </div>
                      </div>
                      {renderTaskDetails(task)}
                    </div>
                  ))}
                </div>
              </Card>
            )}

            {/* Today Tasks */}
            <Card className="p-4 flex-1 flex flex-col">
              <h3 className="font-label-md text-xs font-bold text-primary flex items-center gap-2 mb-3 uppercase tracking-wider">
                <span className="material-symbols-outlined text-[16px]">today</span>
                {t("tasks.today")}
              </h3>
              <div className="space-y-2 flex-1 overflow-y-auto pr-1">
                {todayTasks.map(task => (
                  <div
                    key={task.id}
                    onClick={() => toggleExpand(task.id)}
                    className={`rounded-lg p-3 border transition-colors flex flex-col gap-1 cursor-pointer select-none group ${
                      task.completed
                        ? "border-outline-variant/35 bg-surface/50 opacity-60"
                        : "border-transparent bg-surface-container-low hover:border-primary"
                    }`}
                  >
                    <div className="flex gap-3 w-full">
                      <div onClick={(e) => e.stopPropagation()}>
                        <Checkbox checked={task.completed} onChange={() => void toggleTask(task.id)} className="mt-1" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex justify-between items-start gap-2">
                          <h4 className={`font-label-md text-sm text-on-surface font-semibold group-hover:text-primary transition-colors truncate ${
                            task.completed ? "line-through opacity-50" : ""
                          }`}>
                            {task.title}
                          </h4>
                          <div className="flex items-center gap-1 text-on-surface-variant text-[10px] shrink-0 mt-0.5">
                            <span className="material-symbols-outlined text-[12px]">schedule</span>
                            <span>{task.time}</span>
                          </div>
                        </div>
                        <div className="flex justify-between items-center mt-2">
                          <p className="text-xs text-on-surface-variant truncate mr-2">{task.company}</p>
                          <span className="px-2.5 py-0.5 bg-primary-fixed text-on-primary-fixed-variant text-[10px] font-semibold rounded-full uppercase shrink-0">
                            {task.type}
                          </span>
                        </div>
                      </div>
                    </div>
                    {renderTaskDetails(task)}
                  </div>
                ))}
              </div>
            </Card>

            {/* Upcoming Tasks */}
            <div className="bg-surface-container rounded-xl p-4 border border-outline-variant/50">
              <h3 className="font-label-md text-xs font-bold text-on-surface-variant flex items-center gap-2 mb-3 uppercase tracking-wider">
                <span className="material-symbols-outlined text-[16px]">upcoming</span>
                {t("tasks.upcoming")}
              </h3>
              <div className="space-y-2">
                {upcomingTasks.map(task => (
                  <div 
                    key={task.id} 
                    onClick={() => toggleExpand(task.id)}
                    className="bg-surface rounded-lg p-3 border border-transparent flex flex-col gap-1 cursor-pointer select-none hover:border-outline-variant/50 transition-all"
                  >
                    <div className="flex gap-3 w-full">
                      <div onClick={(e) => e.stopPropagation()}>
                        <Checkbox checked={task.completed} onChange={() => void toggleTask(task.id)} className="mt-1" />
                      </div>
                      <div className="flex-1">
                        <h4 className={`font-label-md text-sm text-on-surface font-semibold ${task.completed ? "line-through opacity-50" : ""}`}>
                          {task.title}
                        </h4>
                        <p className="text-xs text-on-surface-variant mt-0.5">{task.company}</p>
                        <div className="flex items-center gap-1 mt-2 text-on-surface-variant text-[10px]">
                          <span className="material-symbols-outlined text-[12px]">calendar_today</span>
                          <span>{task.time}</span>
                        </div>
                      </div>
                    </div>
                    {renderTaskDetails(task)}
                  </div>
                ))}
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
