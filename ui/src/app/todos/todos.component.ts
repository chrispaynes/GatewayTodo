import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { from, Observable } from 'rxjs';
import {
  Todo,
  StatusFilter,
  DefaultTemplateTitle,
  DefaultTemplateDescription,
} from '../todo';
import { TodoService } from '../todo.service';
import { tap, filter } from 'rxjs/operators';

@Component({
  selector: 'app-todos',
  templateUrl: './todos.component.html',
  styleUrls: ['./todos.component.scss'],
})
export class TodosComponent implements OnInit {
  todos: Todo[] = [];
  todoCount: number = 0; // includes all todos except for the blank template

  statusFilter: StatusFilter = '';
  selectedTodos: number[] = [];

  constructor(private router: Router, private todoService: TodoService) {}

  ngOnInit(): void {
    // use the router URL to determine which todos to display
    this.statusFilter = this.mapRouteToFilter(this.router.url.replace('/', ''));

    this.loadTodos(this.todoService.getAllTodos());
  }

  // loadTodos loads todos from a datasource
  private loadTodos(dataSource: Observable<Todo>) {
    this.todos = [];
    this.selectedTodos = [];
    this.todoCount = 0;

    // add a blank todo template for users create a new todo
    if (['', 'add', 'todo'].includes(this.statusFilter)) {
      this.todos.push(
        new Todo(DefaultTemplateTitle, DefaultTemplateDescription)
      );
    }

    from(dataSource)
      .pipe(
        filter((todo: Todo) => {
          // apply a filter to the todo status when there's a filter
          // otherwise return every todo
          return this.statusFilter == ''
            ? !!todo
            : todo.status == this.statusFilter;
        }),
        tap((t: Todo) => {
          this.todos.push(t);
          this.todoCount++;
        })
      )
      .subscribe();
  }

  // mapRouteToFilter maps the current route to a StatusFilter to
  // filter out Todos by their status
  private mapRouteToFilter(route: string): StatusFilter {
    const filterMap = new Map();

    filterMap.set('progress', 'todo');
    filterMap.set('completed', 'completed');
    filterMap.set('archived', 'archived');
    filterMap.set('new', 'new');

    // if we don't find a filter in the map, don't add a filter
    return !filterMap.get(route) ? '' : filterMap.get(route);
  }

  public toggle() {
    // deselect todos if we have some selected
    if (this.selectedTodos.length > 0) {
      this.selectedTodos = [];
    } else {
      // select all todos with a valid id (we're excluding the blank "template" todo)
      this.selectedTodos = this.todos
        .map((t: Todo) => t.id)
        .filter((id: number) => id > 0);
    }

    // toggle each todo's selected state, except for the blank "template" todo
    this.todos.forEach((todo: Todo) => {
      if (todo.id > 0) {
        todo._isSelected = !todo._isSelected;
      }
    });
  }

  public deleteAll(todoIDs: number[]) {
    this.todoService
      .deleteTodos(todoIDs)
      .pipe(
        tap((wasSuccessful: boolean) => {
          if (wasSuccessful) {
            this.selectedTodos = [];
            this.todos = this.todos.filter(
              (todo: Todo) => !todoIDs.includes(todo.id)
            );
          }
        })
      )
      .subscribe();
  }

  // saveAll saves all updates to todos except for
  // the blank template todo in the first position
  // in the UI's todo grid
  public saveAll(todoIDs: number[]) {
    if (todoIDs.length === 0) {
      return;
    }

    // filter out any todo's that were not in the selection list
    const todos: Todo[] = this.todos.filter((t: Todo) =>
      todoIDs.includes(t.id)
    );
    const updatedTodos = this.todoService.updateTodos({ todos: todos });

    this.loadTodos(updatedTodos);
  }
}
