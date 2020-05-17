import { Component, OnInit, SimpleChanges } from '@angular/core';
import { Router, ActivatedRoute, ParamMap, UrlSegment } from '@angular/router';
import { from } from 'rxjs';
import { Todo, StatusFilter } from '../todo';
import { TodoService } from '../todo.service';
import { tap, filter } from 'rxjs/operators';

@Component({
    selector: 'app-todos',
    templateUrl: './todos.component.html',
    styleUrls: ['./todos.component.scss'],
})
export class TodosComponent implements OnInit {
    todos: Todo[] = [];

    statusFilter: StatusFilter = '';
    selectedTodos: number[] = [];

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private todoService: TodoService
    ) {}

    ngOnInit(): void {
        this.statusFilter = this.mapRouteToFilter(
            this.router.url.replace('/', '')
        );

        console.log(this);
        from(this.todoService.getAllTodos())
            .pipe(
                filter((todo: Todo) => {

                  if (['', 'add', 'todo'].includes(this.statusFilter)) {
                    this.todos.push(new Todo('new title', 'new description'));
                  }

                    // apply a filter to the todo status when there's a filter
                    // otherwise return every todo
                    return this.statusFilter == ''
                        ? !!todo
                        : todo.status == this.statusFilter;
                }),
                tap((x: Todo) => this.todos.push(x))
            )
            .subscribe();
    }

    // mapRouteToFilter maps the current route to a StatusFilter to
    // filter out Todos by their status
    private mapRouteToFilter(route: string): StatusFilter {
        console.log(route);
        const filterMap = new Map();
        filterMap.set('progress', 'todo');
        filterMap.set('completed', 'completed');
        filterMap.set('archived', 'archived');
        filterMap.set('new', 'new');

        return !filterMap.get(route) ? '' : filterMap.get(route);
    }

    public toggle() {
        if (this.selectedTodos.length > 0) {
            this.selectedTodos = [];
        } else {
            // select all todos with a valid id
            this.selectedTodos = this.todos
                .map((t: Todo) => t.id)
                .filter((id: number) => id > 0);
        }

        // toggle each Todo's selected state, except for the blank "template" todo
        this.todos.forEach((todo: Todo) => {
            if (todo.id > 0) {
                todo._isSelected = !todo._isSelected;
            }
        });
    }

    public selectTodo(todo) {
        console.log('selected', todo);
    }

    public deleteAll(todoIDs: number[]) {
        this.todoService
            .deleteTodos(todoIDs)
            .pipe(
                tap((wasSuccessful: boolean) => {
                    if (wasSuccessful) {
                        console.log("success, we're deleint");
                        this.todos = this.todos.filter(
                            (todo: Todo) => !todoIDs.includes(todo.id)
                        );
                    }
                })
            )
            .subscribe();
    }
}
