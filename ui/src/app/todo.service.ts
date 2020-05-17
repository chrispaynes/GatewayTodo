import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of, from } from 'rxjs';
import { map, tap , switchMap } from 'rxjs/operators';
import { Todo, Todos } from './todo';

@Injectable({
    providedIn: 'root',
})
export class TodoService {
    constructor(private _http: HttpClient) {}

    public getAllTodos(): Observable<any> {
        return this._http.get('http://localhost:3000/v1/todos').pipe(
          switchMap((x:Todos) => x.todos),
            // // tap((x) => console.log(x)),
            // map((x: Todos) => x.todos)
        );
    }

    public addTodo(todo: Todo): Observable<Todo[]> {
        return this._http.post('http://localhost:3000/v1/todo', {}).pipe(
            tap((x) => console.log(x)),
            map((x: Todos) => x.todos)
        );
    }

    public updateTodo(todo: Todo): Observable<Todo[]> {
        return this._http.put('http://localhost:3000/v1/todo', {}).pipe(
            tap((x) => console.log(x)),
            map((x: Todos) => x.todos)
        );
    }

    public updateTodos(todos: Todo[]): Observable<Todo[]> {
      return this._http.put('http://localhost:3000/v1/todo', {}).pipe(
          tap((x) => console.log(x)),
          map((x: Todos) => x.todos)
      );
  }

    public deleteTodo(): Observable<any> {
        return this._http.delete('http://localhost:3000/v1/todos').pipe(
            tap((x) => console.log(x)),
            map((x: Todos) => x.todos)
        );
    }

    public deleteTodos(ids: number[]): Observable<any> {
        return this._http.delete('http://localhost:3000/v1/todos').pipe(
            tap((x) => console.log(x)),
            map((x: Todos) => x.todos)
        );
    }

    // getEmail gets emailed from LocalStorage or the mock-email datastore
    getEmail(id: string): Observable<any> {
        // match unbracketed server generated Message-IDs such as:
        // 6426946.1413.1301675117949.JavaMail.tomcat@osadmin02
        if (isNaN(+id)) {
            const localStorageEmail = localStorage.getItem(id);

            // if (localStorageEmail === null) {
            //     return of(
            //         EMAILS.find((email) => email.MessageId === `<${id}>`)
            //     );
            // }

            return of(JSON.parse(localStorageEmail));
        }
    }

    // postEmail POSTs an email message to the server
    postEmail(email: string): Observable<any> {
        const httpOptions = {
            headers: new HttpHeaders({
                'Content-Type': 'text/plain',
                Accept: 'application/json',
            }),
        };

        return this._http.post<any>(
            'http://api-gmc.localhost/email',
            email,
            httpOptions
        );
    }
}

// {
//   "todos": [
//       {
//           "id": "100008",
//           "title": "<string>",
//           "description": "<string>",
//           "status": "new",
//           "createdDT": "2020-05-16T23:01:09.801977Z",
//           "updatedDT": "2020-05-17T01:01:09.789325Z"
//       },
//       {
