import { Injectable } from '@angular/core';
import {
  HttpClient,
  HttpErrorResponse,
  HttpResponse,
} from '@angular/common/http';
import { Observable, of, from } from 'rxjs';
import { map, switchMap, catchError, timeout } from 'rxjs/operators';
import { Todo, Todos } from './todo';

@Injectable({
  providedIn: 'root',
})
export class TodoService {
  constructor(private _http: HttpClient) {}

  public getAllTodos(): Observable<Todo> {
    return this._http.get('http://localhost:3000/v1/todos').pipe(
      timeout(2000),
      catchError((error: HttpErrorResponse) => {
        return of(null);
      }),
      switchMap((t: Todos) => {
        return t.todos ? t.todos : of(null);
      })
    );
  }

  public addTodo(todo: Todo): Observable<Todo> {
    const body = { title: todo.title, description: todo.description };

    return this._http.post('http://localhost:3000/v1/todo', body).pipe(
      timeout(2000),
      catchError((error: HttpErrorResponse) => {
        return of(null);
      }),
      map((t: Todo) => {
        return t ? t : null;
      })
    );
  }

  public updateTodo(todo: Todo, newStatus?: 'completed'): Observable<Todo> {
    const body = {
      id: todo.id,
      title: todo.title,
      description: todo.description,
      status: todo.status,
    };

    return this._http
      .put(`http://localhost:3000/v1/todo/${todo.id}`, body)
      .pipe(
        timeout(2000),
        catchError((error: HttpErrorResponse) => {
          return of(null);
        }),
        map((t: Todo) => {
          return t ? t : null;
        })
      );
  }

  public updateTodos(todos: Todos): Observable<Todo> {
    return this._http.put('http://localhost:3000/v1/todos', todos).pipe(
      timeout(2000),
      catchError((error: HttpErrorResponse) => {
        return of(null);
      }),
      switchMap((t: Todos) => {
        return t.todos ? t.todos : of(null);
      })
    );
  }

  public deleteTodo(id: number): Observable<boolean> {
    return this._http
      .request('delete', `http://localhost:3000/v1/todo/${id}`, {
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
