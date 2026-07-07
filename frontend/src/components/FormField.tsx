interface FormFieldProps {
  label: string
  type: string
  value: string
  onChange: (value: string) => void
  required?: boolean
  minLength?: number
  autoComplete?: string
}

export function FormField({ label, type, value, onChange, required, minLength, autoComplete }: FormFieldProps) {
  return (
    <label className="block">
      <span className="mb-1 block text-sm font-medium text-slate-600 dark:text-slate-300">{label}</span>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        required={required}
        minLength={minLength}
        autoComplete={autoComplete}
        className="w-full rounded-lg border border-slate-300 bg-transparent px-3 py-2 text-sm text-slate-800 outline-none focus:border-indigo-500 dark:border-slate-600 dark:text-slate-100"
      />
    </label>
  )
}
