export class Todo {
    readonly id: number;
    title: string;
    description: string;
    status: string;
    createdDT: string;
    updatedDT: string;

    constructor(
        id: number,
        title: string,
        description: string,
        status: string,
        createdDT: string,
        updatedDT: string
    ) {
        this.id = id;
        this.title = title;
        this.description = description;
        this.status = status;
        this.createdDT = createdDT;
        this.updatedDT = updatedDT;
    }
}

export interface Todos {
  todos: Todo[];
}

export type StatusFilter = '' | 'new' | 'todo' | 'completed' | 'archived';
