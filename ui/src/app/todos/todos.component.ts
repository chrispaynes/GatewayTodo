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

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private todoService: TodoService
    ) {}

    ngOnInit(): void {
        this.statusFilter = this.mapRouteToFilter(this.router.url.replace('/', ''));

        from(this.todoService.getAllTodos())
            .pipe(
                filter((todo: Todo) => {
                    // apply a filter to the todo status when there's a filter
                    // otherwise return every todo
                    return this.statusFilter !== ''
                        ? todo.status == this.statusFilter
                        : !!todo;
                }),
                tap((x: Todo) => this.todos.push(x))
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

        return !filterMap.get(route) ? '' : filterMap.get(route);
    }

    public toggle() {
        this.todos.forEach((todo: Todo) => {});
    }
}
