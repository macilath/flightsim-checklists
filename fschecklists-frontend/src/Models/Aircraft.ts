import { LiteChecklist } from "./Checklist";

export default interface Aircraft {
    id: number
    name: string
    alias: string
    checklists: LiteChecklist[]
}
