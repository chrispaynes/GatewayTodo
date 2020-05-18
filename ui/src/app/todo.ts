export class Todo {
    readonly id: number;
    title: string;
    description: string;
    status: string;
    createdDT: string;
    updatedDT: string;
    _isDirty: boolean;
    _isSelected: boolean;

    constructor(
        title: string,
        description: string,
        id?: number,
        status?: string,
        createdDT?: string,
        updatedDT?: string
    ) {
        this.id = id;
        this.title = title;
        this.description = description;
        this.status = status;
        this.createdDT = createdDT;
        this.updatedDT = updatedDT;
        this._isDirty = false;
        this._isSelected = false;
    }
}

export interface Todos {
  todos: Todo[];
}

export type StatusFilter = '' | 'new' | 'todo' | 'completed' | 'archived';

export type TodoEdit = {"title": string , "description": string}

export const DefaultTemplateTitle = "new title";
export const DefaultTemplateDescription = "new description";
