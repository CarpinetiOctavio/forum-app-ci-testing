# CI/CD Pipeline - GitHub Actions

## ¿Qué es CI/CD?

**CI/CD = Continuous Integration / Continuous Deployment**

- **CI (Integración Continua)**: Los tests se ejecutan automáticamente cada vez que haces push
- **CD (Despliegue Continuo)**: El código se construye/despliega automáticamente si los tests pasan

## Cómo funciona este pipeline

### Trigger (Disparador)
```yaml
on:
  push:
    branches: [ main, master, develop ]
  pull_request:
    branches: [ main, master, develop ]
```

Se ejecuta automáticamente cuando:
- Hacés push a las ramas main/master/develop
- Creás un pull request

### Jobs (Trabajos)

#### 1. Backend Tests (Go)
```bash
go test ./... -v -coverprofile=coverage.out
```
- Ejecuta todos los tests del backend
- Genera reporte de cobertura
- Si falla, detiene el pipeline

#### 2. Frontend Tests (React)
```bash
npm test -- --coverage --watchAll=false
```
- Ejecuta todos los tests del frontend (sin modo watch)
- Genera reporte de cobertura
- Si falla, detiene el pipeline

#### 3. Backend Build
- Depende de: `backend-tests` (solo se ejecuta si los tests pasan)
- Compila el backend con `go build`
- Verifica que el código sea construible

#### 4. Frontend Build
- Depende de: `frontend-tests`
- Compila el frontend con `npm run build`
- Verifica que React compile correctamente

#### 5. Summary
- Se ejecuta siempre (`if: always()`)
- Muestra resumen final del pipeline

### Dependencias

```
backend-tests ──┬──→ backend-build ──┐
                                      └──→ summary
frontend-tests ─┬──→ frontend-build ─┘
```

## Beneficios

1. **Automatización**: No necesitás ejecutar tests manualmente
2. **Consistencia**: Los tests se ejecutan en un servidor limpio (no tu máquina)
3. **Visibilidad**: Todos ven si el código está "roto" o "funciona"
4. **Prevención**: Previene mergear código con tests que fallan
5. **Documentación**: El historio de builds es un registro

## Flujo típico

```
1. Escribís código localmente
2. Ejecutás tests locales (npm test / go test)
3. Haces commit y push
4. GitHub Actions se activa automáticamente
5. Los 5 jobs se ejecutan en paralelo
6. Si todo pasa ✅: Tu PR puede ser mergeado
7. Si algo falla ❌: Ves el error y arreglas
```

## Ver resultados en GitHub

1. Entra a https://github.com/tu-usuario/tp06-testing
2. Click en pestaña "Actions"
3. Verás un listado de workflows ejecutados
4. Click en uno para ver detalles
5. Expande cada job para ver los logs

## Para esta materia (TP06)

Este pipeline demuestra:
- ✅ **Automatización**: Tests se ejecutan sin intervención manual
- ✅ **Validación**: El código debe pasar tests para ser construido
- ✅ **Reproducibilidad**: Se ejecuta en un servidor limpio (como en producción)
- ✅ **DevOps**: Implementación práctica de CI/CD