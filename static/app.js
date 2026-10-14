(() => {
    'use strict';

    const expressionEl = document.getElementById('expression');
    const resultEl = document.getElementById('result');
    const historyEl = document.getElementById('history');

    let currentExpression = '';
    let lastResult = '';
    const history = [];
    const MAX_HISTORY = 5;

    // --- Display ---
    function updateDisplay() {
        expressionEl.textContent = currentExpression || '';
        if (!currentExpression) {
            resultEl.textContent = lastResult || '0';
            resultEl.classList.remove('error');
        }
    }

    function showResult(value) {
        resultEl.textContent = value;
        resultEl.classList.remove('error');
    }

    function showError(message) {
        resultEl.textContent = message;
        resultEl.classList.add('error');
    }

    function addToHistory(expr, result) {
        history.push({ expr, result });
        if (history.length > MAX_HISTORY) {
            history.shift();
        }
        renderHistory();
    }

    function renderHistory() {
        historyEl.innerHTML = '';
        for (const entry of history) {
            const div = document.createElement('div');
            div.className = 'history-entry';
            div.textContent = `${entry.expr} = ${entry.result}`;
            historyEl.appendChild(div);
        }
        historyEl.scrollTop = historyEl.scrollHeight;
    }

    // --- Input ---
    function appendValue(value) {
        currentExpression += value;
        updateDisplay();
    }

    function clearAll() {
        currentExpression = '';
        lastResult = '';
        resultEl.textContent = '0';
        resultEl.classList.remove('error');
        updateDisplay();
    }

    function backspace() {
        currentExpression = currentExpression.slice(0, -1);
        updateDisplay();
    }

    // --- Calculation ---
    async function calculate() {
        const expr = currentExpression.trim();
        if (!expr) return;

        // Replace display characters with actual operators
        const sanitized = expr
            .replace(/×/g, '*')
            .replace(/÷/g, '/')
            .replace(/−/g, '-');

        try {
            const response = await fetch('/api/calculate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ expression: sanitized }),
            });

            const data = await response.json();

            if (data.error) {
                showError(data.error);
            } else {
                const formatted = formatResult(data.result);
                showResult(formatted);
                addToHistory(expr, formatted);
                lastResult = formatted;
                currentExpression = '';
                updateDisplay();
            }
        } catch (err) {
            showError('Verbindungsfehler');
        }
    }

    function formatResult(num) {
        // Avoid floating-point display artifacts
        if (Number.isInteger(num)) {
            return num.toString();
        }
        // Up to 10 significant decimals, strip trailing zeros
        return parseFloat(num.toPrecision(12)).toString();
    }

    // --- Button Clicks ---
    document.querySelector('.keypad').addEventListener('click', (e) => {
        const btn = e.target.closest('.btn');
        if (!btn) return;

        const value = btn.dataset.value;
        const action = btn.dataset.action;

        if (action === 'clear') {
            clearAll();
        } else if (action === 'backspace') {
            backspace();
        } else if (action === 'equals') {
            calculate();
        } else if (value) {
            appendValue(value);
        }
    });

    // --- Keyboard Support ---
    document.addEventListener('keydown', (e) => {
        const key = e.key;

        // Prevent default for calculator keys
        if (['Enter', 'Escape', 'Backspace'].includes(key) ||
            /^[0-9+\-*/.()=]$/.test(key)) {
            e.preventDefault();
        }

        if (key >= '0' && key <= '9') {
            appendValue(key);
        } else if (key === '.') {
            appendValue('.');
        } else if (key === '+') {
            appendValue('+');
        } else if (key === '-') {
            appendValue('-');
        } else if (key === '*') {
            appendValue('*');
        } else if (key === '/') {
            appendValue('/');
        } else if (key === '(') {
            appendValue('(');
        } else if (key === ')') {
            appendValue(')');
        } else if (key === 'Enter' || key === '=') {
            calculate();
        } else if (key === 'Backspace') {
            backspace();
        } else if (key === 'Escape') {
            clearAll();
        }
    });

    // Initial display
    updateDisplay();
})();

