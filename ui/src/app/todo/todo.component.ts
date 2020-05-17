import { Component, OnInit, Input, HostListener } from '@angular/core';
import {
    FormControl,
    FormArray,
    FormGroup,
    FormBuilder,
    Validators,
} from '@angular/forms';

import { TodoService } from '../todo.service';
import { Todo, TodoEdit } from '../todo';
import { map, tap } from 'rxjs/operators';

@Component({
    selector: 'app-todo',
    templateUrl: './todo.component.html',
    styleUrls: ['./todo.component.scss'],
})
export class TodoComponent implements OnInit {
    @Input() todo: Todo;

    @HostListener('document:keydown.escape', ['$event'])
    onKeydownHandler(event: KeyboardEvent) {
        console.log(event);
        this.editMode = false;
    }

    controls: FormArray;
    editMode: boolean = false;
    originalTitle: string = '';
    originalDescription: string = '';
    form: FormGroup;

    constructor(
        private _todoService: TodoService,
        private formBuilder: FormBuilder
    ) {}

    ngOnInit(): void {
        this.form = this.formBuilder.group({
            title: [this.todo.title, Validators.required],
            description: [this.todo.description, Validators.required],
        });
        this.originalTitle = this.todo.title;
        this.originalDescription = this.todo.description;

        this.onChanges();
    }

    onChanges(): void {
        this.form.valueChanges.subscribe((edit: TodoEdit) => {
            if (
                edit.title !== this.originalTitle ||
                edit.description !== this.originalDescription
            ) {
                Object.assign(this.todo, { ...edit, _isDirty: true });
            } else {
                Object.assign(this.todo, { ...edit, _isDirty: false });
            }

            console.log('val', this.todo);
        });
    }

    public onFocus(event) {
        if (!event) {
            this.editMode = false;
        }

        console.log('FOCUS', event);
    }

    public onEdit(todo: Todo) {
        this.editMode = true;
        console.log('onEdit', todo);
    }

    public onDelete(todo: Todo) {
        console.log('onDelete', todo);
    }

    public onAdd(todo: Todo) {
        console.log('onAdd', todo);
    }

    public onUpdate(todo: Todo) {
        console.log('onUpdate', todo);
    }

    public onToggleCompletion(todo: Todo) {
        console.log('onToggleCompletion', todo);
    }

    public onSelect(event) {
      // console.log("selected before", event);
      // this.todo._isSelected = !event._isSelected
      // console.log("selected after", this.todo);
    }

}
