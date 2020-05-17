import { Component, OnInit, Input } from '@angular/core';
import { FormControl } from '@angular/forms';

import { TodoService } from '../todo.service';
import { Todo } from '../todo';
import { map, tap } from 'rxjs/operators';

@Component({
    selector: 'app-todo',
    templateUrl: './todo.component.html',
    styleUrls: ['./todo.component.scss'],
})
export class TodoComponent implements OnInit {
    @Input() todo: Todo;

    constructor(private _todoService: TodoService) {}

    ngOnInit(): void {
      if (!this.todo) {
        // create a template for creating a new todo
        this.todo = new Todo(0, "", "", "", "", "")
        console.log("looks like we need to create a todo");
      }
    }


    title = new FormControl('');

    updateName() {
      this.title.setValue('Nancy');
    }

    public onEdit(todo: Todo) {
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

    // emails: Email[];
    // selectedEmail: Email;

    // onSelect(email: Email): void {
    //   this.selectedEmail = email;
    // }
}
