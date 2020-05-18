import { Component, OnInit, Input, HostListener } from '@angular/core';
import {
  FormArray,
  FormGroup,
  FormBuilder,
  Validators,
} from '@angular/forms';

import { TodoService } from '../todo.service';
import {
  Todo,
  TodoEdit,
  DefaultTemplateDescription,
  DefaultTemplateTitle,
} from '../todo';
import { tap } from 'rxjs/operators';
import { from, Observable } from 'rxjs';

@Component({
  selector: 'app-todo',
  templateUrl: './todo.component.html',
  styleUrls: ['./todo.component.scss'],
})
export class TodoComponent implements OnInit {
  @Input() todo: Todo;

  @HostListener('document:keydown.escape', ['$event'])
  onKeydownHandler(event: KeyboardEvent) {
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
      // if the todo's updated title or description don't match its original,
      // assign the properties and consider the todo dirty. otherwise consider
      // the todo unchangeds
      if (
        edit.title !== this.originalTitle ||
        edit.description !== this.originalDescription
      ) {
        Object.assign(this.todo, {
          ...edit,
          _isDirty: true,
        });
      } else {
        Object.assign(this.todo, {
          ...edit,
          _isDirty: false,
        });
      }
    });
  }

  public onFocus(event) {
    if (!event) {
      this.editMode = false;
    }
  }

  public onEdit(todo: Todo) {
    this.editMode = true;
  }

  public onDelete(todoID: number) {
    this._todoService
      .deleteTodo(todoID)
      .pipe(
        tap((wasSuccessful: boolean) => {
          if (wasSuccessful) {
            location.reload();
          }
        })
      )
      .subscribe();
  }

  public onAdd(todo: Todo) {
    if (this.shouldDisableSave(todo)) {
      return;
    }

    this.reloadTodo(this._todoService.addTodo(todo));
  }

  public shouldDisableSave(todo: Todo): boolean {
    if (!todo.id) {
      return (
        todo.title === '' ||
        todo.description == '' ||
        todo.title === DefaultTemplateTitle ||
        todo.description == DefaultTemplateDescription
      );
    }

    return !todo._isDirty || todo.title === '' || todo.description == '';
  }

  public shouldEnableCompleteButton(todo:Todo): boolean {
    return !this.editMode && todo.id && !todo._isDirty && todo.status !== 'completed'
  }

  public onUpdate(todo: Todo, newStatus?: 'completed') {
    if (newStatus) {
      todo.status = newStatus;
    }

    if (!newStatus && this.shouldDisableSave(todo)) {
      return;
    }

    this.reloadTodo(this._todoService.updateTodo(todo));
  }

  // loadTodos loads todos from a datasource
  private reloadTodo(dataSource: Observable<Todo>) {
    from(dataSource)
      .pipe(
        tap((todo: Todo) => {
          if (todo) {
            this.todo = todo;
            this.editMode = false;
          }
        })
      )
      .subscribe();
  }
}
