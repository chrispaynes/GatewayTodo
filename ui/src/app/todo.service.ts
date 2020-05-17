import { Injectable } from '@angular/core';
import {
    HttpClient,
    HttpHeaders,
    HttpParams,
    HttpErrorResponse,
    HttpResponse,
} from '@angular/common/http';
import { Observable, of, from } from 'rxjs';
import { map, tap, switchMap, catchError, timeout } from 'rxjs/operators';
import { Todo, Todos } from './todo';

@Injectable({
    providedIn: 'root',
})
export class TodoService {
    constructor(private _http: HttpClient) {}
    todos$: any = [];

    public getAllTodos(): Observable<Todo> {
        return this._http.get('http://localhost:3000/v1/todos').pipe(
            timeout(2000),
            catchError((error: HttpErrorResponse) => {
                return of(null);
            }),
            switchMap((x: Todos) => {
              this.todos$ = x;
              return x.todos
            })
        );
    }

    public addTodo(todo: Todo): Observable<Todo[]> {
        return this._http.post('http://localhost:3000/v1/todo', {}).pipe(
            timeout(2000),
            catchError((error: HttpErrorResponse) => {
                return of(null);
            }),
            tap((x) => console.log(x)),
            map((x: Todos) => x.todos)
        );
    }

    public updateTodo(todo: Todo): Observable<Todo[]> {
        return this._http.put('http://localhost:3000/v1/todo', {}).pipe(
            timeout(2000),
            catchError((error: HttpErrorResponse) => {
                return of(null);
            }),
            tap((x) => console.log(x)),
            map((x: Todos) => x.todos)
        );
    }

    public updateTodos(todos: Todo[]): Observable<Todo[]> {
        return this._http.put('http://localhost:3000/v1/todo', {}).pipe(
            timeout(2000),
            catchError((error: HttpErrorResponse) => {
                return of(null);
            }),
            tap((x) => console.log(x)),
            map((x: Todos) => x.todos)
        );
    }

    public deleteTodo(id: number): Observable<boolean> {
        return this._http
            .request('delete', `http://localhost:3000/v1/todo${id}`, {
                observe: 'response',
            })
            .pipe(
                timeout(2000),
                catchError((error: HttpErrorResponse) => {
                    return of(false);
                }),
                map((resp: HttpResponse<Object>) => resp && resp.status == 200)
            );
    }

    public deleteTodos(ids: number[]): Observable<boolean> {
        return this._http
            .request('delete', 'http://localhost:3000/v1/todos', {
                observe: 'response',
                body: { ids: ids },
            })
            .pipe(
                timeout(2000),
                catchError((error: HttpErrorResponse) => {
                    return of(false);
                }),
                map((resp: HttpResponse<Object>) => resp && resp.status == 200)
            );
    }
}
