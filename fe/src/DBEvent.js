import React, { useContext, useState } from 'react';
import { Form, Button, Spinner, Alert, Row, Col, Tabs, Tab, Collapse } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import { ThemeContext } from "./ThemeContext";
import ObjectLinkSelector from './ObjectLinkSelector';
import { ObjectSearch } from "./DBObject";
import PermissionsEditor from './PermissionsEditor';
import { formatDescription } from './sitenavigation_utils';


/**
 * EventView - Display component for DBEvent objects
 * Shows event details including dates, times, category, and optional URL
 */
export function EventView({ data, metadata, objectData, dark }) {
    const { t } = useTranslation();

    const isDeleted = data && data.deleted_date;

    // Parse dates
    const startDate = data.start_date ? new Date(data.start_date) : null;
    const endDate = data.end_date ? new Date(data.end_date) : null;
    const isAllDay = data.all_day === '1' || data.all_day === 1;

    // Format date and time
    const formatDateTime = (date) => {
        if (!date) return '';
        const dateStr = date.toLocaleDateString();
        const timeStr = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        return { dateStr, timeStr };
    };

    const startFormatted = startDate ? formatDateTime(startDate) : null;
    const endFormatted = endDate ? formatDateTime(endDate) : null;

    return (
        <div style={isDeleted ? { opacity: 0.5 } : {}}>
            <h2 className={dark ? 'text-light' : 'text-dark'}>
                <i className="bi bi-calendar-event me-2"></i>
                {data.name}{isDeleted ? ' ('+t('dbobjects.deleted')+')' : ''}
            </h2>

            {data.description && (
                <>
                    <hr />
                    <div dangerouslySetInnerHTML={{ __html: formatDescription(data.description) }}></div>
                </>
            )}

            <hr />

            <div className="event-details">
                <Row className="mb-3">
                    <Col md={6}>
                        <div className="mb-2">
                            <strong><i className="bi bi-calendar-check me-2"></i>{t('event.start') || 'Start'}:</strong>
                            <div className="ms-4">
                                {startFormatted && (
                                    <>
                                        <div>{startFormatted.dateStr}</div>
                                        {!isAllDay && <div className="text-muted">{startFormatted.timeStr}</div>}
                                    </>
                                )}
                            </div>
                        </div>
                    </Col>
                    <Col md={6}>
                        <div className="mb-2">
                            <strong><i className="bi bi-calendar-x me-2"></i>{t('event.end') || 'End'}:</strong>
                            <div className="ms-4">
                                {endFormatted && (
                                    <>
                                        <div>{endFormatted.dateStr}</div>
                                        {!isAllDay && <div className="text-muted">{endFormatted.timeStr}</div>}
                                    </>
                                )}
                            </div>
                        </div>
                    </Col>
                </Row>

                {isAllDay && (
                    <div className="mb-2">
                        <span className="badge bg-info">
                            <i className="bi bi-clock me-1"></i>
                            {t('event.all_day') || 'All Day Event'}
                        </span>
                    </div>
                )}

                {data.category && (
                    <div className="mb-2">
                        <strong><i className="bi bi-tag me-2"></i>{t('event.category') || 'Category'}:</strong>
                        <span className="ms-2 badge bg-secondary">{data.category}</span>
                    </div>
                )}

                {data.url && (
                    <div className="mb-3">
                        <strong><i className="bi bi-link-45deg me-2"></i>{t('event.url') || 'URL'}:</strong>
                        <div className="ms-4">
                            <a href={data.url} target="_blank" rel="noopener noreferrer">
                                {data.url}
                                <i className="bi bi-box-arrow-up-right ms-2"></i>
                            </a>
                        </div>
                    </div>
                )}

                {data.alarm === '1' && (
                    <div className="mb-2">
                        <span className="badge bg-warning text-dark">
                            <i className="bi bi-bell me-1"></i>
                            {t('event.alarm_set') || 'Alarm Set'}
                            {data.alarm_minute && data.alarm_unit !== undefined && (
                                <>
                                    {' - '}
                                    {data.alarm_minute} 
                                    {' '}
                                    {data.alarm_unit === '0' ? (t('event.minutes') || 'minutes') : 
                                     data.alarm_unit === '1' ? (t('event.hours') || 'hours') : 
                                     (t('event.days') || 'days')}
                                    {' '}
                                    {data.before_event === '0' ? (t('event.before_start') || 'before') : (t('event.after_start') || 'after')}
                                </>
                            )}
                        </span>
                    </div>
                )}

                {data.recurrence === '1' && (
                    <div className="mb-3 p-3 border rounded">
                        <strong><i className="bi bi-arrow-repeat me-2"></i>{t('event.recurrence') || 'Recurrence'}:</strong>
                        <div className="ms-4 mt-2">
                            {data.recurrence_type === '0' && (
                                <div>
                                    <i className="bi bi-calendar-day me-2"></i>
                                    {t('event.daily') || 'Daily'} - {t('event.every') || 'Every'} {data.daily_every_x || 1} {t('event.days') || 'day(s)'}
                                </div>
                            )}
                            {data.recurrence_type === '1' && (
                                <div>
                                    <i className="bi bi-calendar-week me-2"></i>
                                    {t('event.weekly') || 'Weekly'} - {t('event.every') || 'Every'} {data.weekly_every_x || 1} {t('event.weeks') || 'week(s)'}
                                    {data.weekly_day_of_the_week !== undefined && (
                                        <>
                                            {' '}{t('event.on') || 'on'}{' '}
                                            {['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'][parseInt(data.weekly_day_of_the_week)] && 
                                             t(`event.${['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'][parseInt(data.weekly_day_of_the_week)]}`)}
                                        </>
                                    )}
                                </div>
                            )}
                            {data.recurrence_type === '2' && (
                                <div>
                                    <i className="bi bi-calendar-month me-2"></i>
                                    {t('event.monthly') || 'Monthly'} - {t('event.every') || 'Every'} {data.monthly_every_x || 1} {t('event.months') || 'month(s)'}
                                    {data.monthly_day_of_the_month > 0 && (
                                        <> {t('event.on_day') || 'on day'} {data.monthly_day_of_the_month}</>
                                    )}
                                    {data.monthly_week_number > 0 && (
                                        <>
                                            {' '}{t('event.on_the') || 'on the'}{' '}
                                            {[t('event.first'), t('event.second'), t('event.third'), t('event.fourth'), t('event.last')][parseInt(data.monthly_week_number) - 1]}
                                            {' '}
                                            {['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'][parseInt(data.monthly_week_day)] && 
                                             t(`event.${['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'][parseInt(data.monthly_week_day)]}`)}
                                        </>
                                    )}
                                </div>
                            )}
                            {data.recurrence_type === '3' && (
                                <div>
                                    <i className="bi bi-calendar-range me-2"></i>
                                    {t('event.yearly') || 'Yearly'}
                                    {data.yearly_day_of_the_year > 0 ? (
                                        <> - {t('event.day') || 'Day'} {data.yearly_day_of_the_year} {t('event.of_year') || 'of the year'}</>
                                    ) : (
                                        <>
                                            {data.yearly_month_number > 0 && data.yearly_month_day > 0 && (
                                                <>
                                                    {' - '}
                                                    {[t('event.january'), t('event.february'), t('event.march'), t('event.april'), 
                                                      t('event.may'), t('event.june'), t('event.july'), t('event.august'),
                                                      t('event.september'), t('event.october'), t('event.november'), t('event.december')][parseInt(data.yearly_month_number) - 1]}
                                                    {' '}{data.yearly_month_day}
                                                </>
                                            )}
                                            {data.yearly_week_number > 0 && (
                                                <>
                                                    {' - '}
                                                    {[t('event.first'), t('event.second'), t('event.third'), t('event.fourth'), t('event.last')][parseInt(data.yearly_week_number) - 1]}
                                                    {' '}
                                                    {['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'][parseInt(data.yearly_week_day)] && 
                                                     t(`event.${['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'][parseInt(data.yearly_week_day)]}`)}
                                                    {' '}{t('event.of') || 'of'}{' '}
                                                    {[t('event.january'), t('event.february'), t('event.march'), t('event.april'), 
                                                      t('event.may'), t('event.june'), t('event.july'), t('event.august'),
                                                      t('event.september'), t('event.october'), t('event.november'), t('event.december')][parseInt(data.yearly_month_number) - 1]}
                                                </>
                                            )}
                                        </>
                                    )}
                                </div>
                            )}
                            
                            {(data.recurrence_times > 0 || data.recurrence_end_date) && (
                                <div className="mt-2 text-muted small">
                                    {data.recurrence_times > 0 && (
                                        <>{t('event.repeat') || 'Repeat'} {data.recurrence_times} {t('event.times') || 'times'}</>
                                    )}
                                    {data.recurrence_end_date && (
                                        <>
                                            {t('event.until') || 'Until'} {new Date(data.recurrence_end_date).toLocaleDateString()}
                                        </>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}


/**
 * EventEdit - Edit form component for DBEvent objects
 * Basic fields for now: name, description, start_date, end_date, all_day, url, category
 * Recurrence fields are available but not yet implemented in UI
 */
export function EventEdit({ data, onSave, onCancel, onDelete, saving, error, dark }) {
    const { t } = useTranslation();
    
    const isDeleted = data && data.deleted_date;

    // Helper to format datetime for input[type="datetime-local"]
    const formatDateTimeLocal = (dateStr) => {
        if (!dateStr) return '';
        const date = new Date(dateStr);
        if (isNaN(date.getTime())) return '';
        // Format: YYYY-MM-DDTHH:mm
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        return `${year}-${month}-${day}T${hours}:${minutes}`;
    };

    const [formData, setFormData] = useState({
        name: data.name || '',
        description: data.description || '',
        start_date: formatDateTimeLocal(data.start_date) || '',
        end_date: formatDateTimeLocal(data.end_date) || '',
        all_day: data.all_day === '1' || data.all_day === 1 ? '1' : '0',
        url: data.url || '',
        category: data.category || '',
        alarm: data.alarm === '1' || data.alarm === 1 ? '1' : '0',
        alarm_minute: data.alarm_minute || '15',
        alarm_unit: data.alarm_unit || '0', // 0=minutes, 1=hours, 2=days
        before_event: data.before_event === '1' || data.before_event === 1 ? '1' : '0',
        recurrence: data.recurrence === '1' || data.recurrence === 1 ? '1' : '0',
        recurrence_type: data.recurrence_type || '0', // 0=Daily, 1=Weekly, 2=Monthly, 3=Yearly
        // Daily
        daily_every_x: data.daily_every_x || 1,
        // Weekly
        weekly_every_x: data.weekly_every_x || 1,
        weekly_day_of_the_week: data.weekly_day_of_the_week || 0,
        // Monthly
        monthly_every_x: data.monthly_every_x || 1,
        monthly_day_of_the_month: data.monthly_day_of_the_month || 0,
        monthly_week_number: data.monthly_week_number || 0,
        monthly_week_day: data.monthly_week_day || 0,
        // Yearly
        yearly_month_number: data.yearly_month_number || 1,
        yearly_month_day: data.yearly_month_day || 1,
        yearly_week_number: data.yearly_week_number || 0,
        yearly_week_day: data.yearly_week_day || 0,
        yearly_day_of_the_year: data.yearly_day_of_the_year || 0,
        // Recurrence range
        recurrence_times: data.recurrence_times || '0', // 0=always
        recurrence_end_date: formatDateTimeLocal(data.recurrence_end_date) || null,
        fk_obj_id: data.fk_obj_id || '0',
        permissions: data.permissions || 'rwxr-----',
        father_id: data.father_id || '0',
        owner: data.owner || null,
        group_id: data.group_id || null,
    });

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? (checked ? '1' : '0') : value
        }));
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        
        // Convert datetime-local format back to MySQL datetime format
        const submitData = { ...formData };
        if (submitData.start_date) {
            submitData.start_date = new Date(submitData.start_date).toISOString().slice(0, 19).replace('T', ' ');
        }
        if (submitData.end_date) {
            submitData.end_date = new Date(submitData.end_date).toISOString().slice(0, 19).replace('T', ' ');
        }
        if (submitData.recurrence_end_date) {
            submitData.recurrence_end_date = new Date(submitData.recurrence_end_date).toISOString().slice(0, 19).replace('T', ' ');
        }
        
        onSave(submitData);
    };

    const weekDays = [
        { value: '0', label: t('event.monday') || 'Monday' },
        { value: '1', label: t('event.tuesday') || 'Tuesday' },
        { value: '2', label: t('event.wednesday') || 'Wednesday' },
        { value: '3', label: t('event.thursday') || 'Thursday' },
        { value: '4', label: t('event.friday') || 'Friday' },
        { value: '5', label: t('event.saturday') || 'Saturday' },
        { value: '6', label: t('event.sunday') || 'Sunday' },
    ];

    const alarmUnits = [
        { value: '0', label: t('event.minutes') || 'Minutes' },
        { value: '1', label: t('event.hours') || 'Hours' },
        { value: '2', label: t('event.days') || 'Days' },
    ];

    return (
        <Form onSubmit={handleSubmit}>

            <div className="row">
                <div className="col-md-4 mb-3">
                <ObjectLinkSelector
                    value={formData.father_id || '0'}
                    onChange={handleChange}
                    classname="DBObject"
                    fieldName="father_id"
                    label={t('dbobjects.parent')}
                />
                </div>
                <div className="col-md-4 mb-3">
                <ObjectLinkSelector
                    value={formData.owner}
                    onChange={handleChange}
                    classname="DBUser"
                    fieldName="owner"
                    label={t('permissions.owner')}
                    required={false}
                />
                </div>
                <div className="col-md-4 mb-3">
                <ObjectLinkSelector
                    value={formData.group_id}
                    onChange={handleChange}
                    classname="DBGroup"
                    fieldName="group_id"
                    label={t('permissions.group')}
                    required={false}
                />
                </div>
            </div>

            <PermissionsEditor
                value={formData.permissions}
                onChange={handleChange}
                name="permissions"
                label={t('permissions.current') || 'Permissions'}
                dark={dark}
            />

            <Form.Group className="mb-3">
                <Form.Label>{t('common.name')}</Form.Label>
                <Form.Control
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('common.description')}</Form.Label>
                <Form.Control
                    as="textarea"
                    name="description"
                    rows={3}
                    value={formData.description}
                    onChange={handleChange}
                />
            </Form.Group>

            <Row>
                <Col md={6}>
                    <Form.Group className="mb-3">
                        <Form.Label>{t('event.start') || 'Start Date & Time'}</Form.Label>
                        <Form.Control
                            type="datetime-local"
                            name="start_date"
                            value={formData.start_date}
                            onChange={handleChange}
                            required
                        />
                    </Form.Group>
                </Col>
                <Col md={6}>
                    <Form.Group className="mb-3">
                        <Form.Label>{t('event.end') || 'End Date & Time'}</Form.Label>
                        <Form.Control
                            type="datetime-local"
                            name="end_date"
                            value={formData.end_date}
                            onChange={handleChange}
                            required
                        />
                    </Form.Group>
                </Col>
            </Row>

            <Form.Group className="mb-3">
                <Form.Check
                    type="checkbox"
                    name="all_day"
                    label={t('event.all_day') || 'All Day Event'}
                    checked={formData.all_day === '1'}
                    onChange={handleChange}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('event.category') || 'Category'}</Form.Label>
                <Form.Control
                    type="text"
                    name="category"
                    value={formData.category}
                    onChange={handleChange}
                    placeholder={t('event.category_placeholder') || 'e.g., Meeting, Birthday, Holiday'}
                />
            </Form.Group>

            <Form.Group className="mb-3">
                <Form.Label>{t('event.url') || 'URL'}</Form.Label>
                <Form.Control
                    type="url"
                    name="url"
                    value={formData.url}
                    onChange={handleChange}
                    placeholder="https://example.com"
                />
            </Form.Group>

            {/* Alarm Section */}
            <div className="mb-4 p-3 border rounded">
                <Form.Group className="mb-3">
                    <Form.Check
                        type="checkbox"
                        name="alarm"
                        label={<strong>{t('event.alarm') || 'Set Alarm'}</strong>}
                        checked={formData.alarm === '1'}
                        onChange={handleChange}
                    />
                </Form.Group>

                <Collapse in={formData.alarm === '1'}>
                    <div>
                        <Row>
                            <Col md={4}>
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('event.alarm_time') || 'Time'}</Form.Label>
                                    <Form.Control
                                        type="number"
                                        name="alarm_minute"
                                        value={formData.alarm_minute}
                                        onChange={handleChange}
                                        min="0"
                                    />
                                </Form.Group>
                            </Col>
                            <Col md={4}>
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('event.alarm_unit') || 'Unit'}</Form.Label>
                                    <Form.Select
                                        name="alarm_unit"
                                        value={formData.alarm_unit}
                                        onChange={handleChange}
                                    >
                                        {alarmUnits.map(unit => (
                                            <option key={unit.value} value={unit.value}>{unit.label}</option>
                                        ))}
                                    </Form.Select>
                                </Form.Group>
                            </Col>
                            <Col md={4}>
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('event.alarm_when') || 'When'}</Form.Label>
                                    <Form.Select
                                        name="before_event"
                                        value={formData.before_event}
                                        onChange={handleChange}
                                    >
                                        <option value="0">{t('event.before_start') || 'Before start'}</option>
                                        <option value="1">{t('event.after_start') || 'After start'}</option>
                                    </Form.Select>
                                </Form.Group>
                            </Col>
                        </Row>
                    </div>
                </Collapse>
            </div>

            {/* Recurrence Section */}
            <div className="mb-4 p-3 border rounded">
                <Form.Group className="mb-3">
                    <Form.Check
                        type="checkbox"
                        name="recurrence"
                        label={<strong>{t('event.recurrence') || 'Recurring Event'}</strong>}
                        checked={formData.recurrence === '1'}
                        onChange={handleChange}
                    />
                </Form.Group>

                <Collapse in={formData.recurrence === '1'}>
                    <div>
                        <Tabs
                            activeKey={formData.recurrence_type}
                            onSelect={(k) => handleChange({ target: { name: 'recurrence_type', value: k } })}
                            className="mb-3"
                        >
                            {/* Daily Tab */}
                            <Tab eventKey="0" title={t('event.daily') || 'Daily'}>
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('event.every_x_days') || 'Every X days'}</Form.Label>
                                    <Form.Control
                                        type="number"
                                        name="daily_every_x"
                                        value={formData.daily_every_x}
                                        onChange={handleChange}
                                        min="0"
                                    />
                                    <Form.Text className="text-muted">
                                        {t('event.daily_help') || 'Event repeats every X days'}
                                    </Form.Text>
                                </Form.Group>
                            </Tab>

                            {/* Weekly Tab */}
                            <Tab eventKey="1" title={t('event.weekly') || 'Weekly'}>
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('event.every_x_weeks') || 'Every X weeks'}</Form.Label>
                                    <Form.Control
                                        type="number"
                                        name="weekly_every_x"
                                        value={formData.weekly_every_x}
                                        onChange={handleChange}
                                        min="0"
                                    />
                                </Form.Group>
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('event.day_of_week') || 'Day of the week'}</Form.Label>
                                    <Form.Select
                                        name="weekly_day_of_the_week"
                                        value={formData.weekly_day_of_the_week}
                                        onChange={handleChange}
                                    >
                                        {weekDays.map(day => (
                                            <option key={day.value} value={day.value}>{day.label}</option>
                                        ))}
                                    </Form.Select>
                                </Form.Group>
                            </Tab>

                            {/* Monthly Tab */}
                            <Tab eventKey="2" title={t('event.monthly') || 'Monthly'}>
                                <Form.Group className="mb-3">
                                    <Form.Label>{t('event.every_x_months') || 'Every X months'}</Form.Label>
                                    <Form.Control
                                        type="number"
                                        name="monthly_every_x"
                                        value={formData.monthly_every_x}
                                        onChange={handleChange}
                                        min="0"
                                    />
                                </Form.Group>
                                
                                <div className="mb-3 p-2 border rounded">
                                    <Form.Label>{t('event.monthly_pattern') || 'Pattern'}</Form.Label>
                                    
                                    <Form.Group className="mb-2">
                                        <Form.Label>{t('event.day_of_month') || 'Day of the month'}</Form.Label>
                                        <Form.Control
                                            type="number"
                                            name="monthly_day_of_the_month"
                                            value={formData.monthly_day_of_the_month}
                                            onChange={handleChange}
                                            min="0"
                                            max="31"
                                        />
                                        <Form.Text className="text-muted">
                                            {t('event.monthly_day_help') || '0 = not used, 1-31 = specific day'}
                                        </Form.Text>
                                    </Form.Group>

                                    <div className="text-center my-2 text-muted">{t('common.or') || 'OR'}</div>

                                    <Row>
                                        <Col md={6}>
                                            <Form.Group className="mb-2">
                                                <Form.Label>{t('event.week_number') || 'Week number'}</Form.Label>
                                                <Form.Select
                                                    name="monthly_week_number"
                                                    value={formData.monthly_week_number}
                                                    onChange={handleChange}
                                                >
                                                    <option value="0">{t('event.not_used') || 'Not used'}</option>
                                                    <option value="1">{t('event.first') || 'First'}</option>
                                                    <option value="2">{t('event.second') || 'Second'}</option>
                                                    <option value="3">{t('event.third') || 'Third'}</option>
                                                    <option value="4">{t('event.fourth') || 'Fourth'}</option>
                                                    <option value="5">{t('event.last') || 'Last'}</option>
                                                </Form.Select>
                                            </Form.Group>
                                        </Col>
                                        <Col md={6}>
                                            <Form.Group className="mb-2">
                                                <Form.Label>{t('event.week_day') || 'Day of week'}</Form.Label>
                                                <Form.Select
                                                    name="monthly_week_day"
                                                    value={formData.monthly_week_day}
                                                    onChange={handleChange}
                                                >
                                                    {weekDays.map(day => (
                                                        <option key={day.value} value={day.value}>{day.label}</option>
                                                    ))}
                                                </Form.Select>
                                            </Form.Group>
                                        </Col>
                                    </Row>
                                </div>
                            </Tab>

                            {/* Yearly Tab */}
                            <Tab eventKey="3" title={t('event.yearly') || 'Yearly'}>
                                <div className="mb-3 p-2 border rounded">
                                    <Form.Label>{t('event.yearly_pattern') || 'Pattern'}</Form.Label>
                                    
                                    <Row>
                                        <Col md={6}>
                                            <Form.Group className="mb-2">
                                                <Form.Label>{t('event.month') || 'Month'}</Form.Label>
                                                <Form.Select
                                                    name="yearly_month_number"
                                                    value={formData.yearly_month_number}
                                                    onChange={handleChange}
                                                >
                                                    <option value="0">{t('event.not_used') || 'Not used'}</option>
                                                    <option value="1">{t('event.january') || 'January'}</option>
                                                    <option value="2">{t('event.february') || 'February'}</option>
                                                    <option value="3">{t('event.march') || 'March'}</option>
                                                    <option value="4">{t('event.april') || 'April'}</option>
                                                    <option value="5">{t('event.may') || 'May'}</option>
                                                    <option value="6">{t('event.june') || 'June'}</option>
                                                    <option value="7">{t('event.july') || 'July'}</option>
                                                    <option value="8">{t('event.august') || 'August'}</option>
                                                    <option value="9">{t('event.september') || 'September'}</option>
                                                    <option value="10">{t('event.october') || 'October'}</option>
                                                    <option value="11">{t('event.november') || 'November'}</option>
                                                    <option value="12">{t('event.december') || 'December'}</option>
                                                </Form.Select>
                                            </Form.Group>
                                        </Col>
                                        <Col md={6}>
                                            <Form.Group className="mb-2">
                                                <Form.Label>{t('event.day') || 'Day'}</Form.Label>
                                                <Form.Control
                                                    type="number"
                                                    name="yearly_month_day"
                                                    value={formData.yearly_month_day}
                                                    onChange={handleChange}
                                                    min="0"
                                                    max="31"
                                                />
                                            </Form.Group>
                                        </Col>
                                    </Row>

                                    <div className="text-center my-2 text-muted">{t('common.or') || 'OR'}</div>

                                    <Row>
                                        <Col md={4}>
                                            <Form.Group className="mb-2">
                                                <Form.Label>{t('event.month') || 'Month'}</Form.Label>
                                                <Form.Select
                                                    name="yearly_month_number"
                                                    value={formData.yearly_month_number}
                                                    onChange={handleChange}
                                                >
                                                    <option value="0">{t('event.not_used') || 'Not used'}</option>
                                                    <option value="1">{t('event.january') || 'January'}</option>
                                                    <option value="2">{t('event.february') || 'February'}</option>
                                                    <option value="3">{t('event.march') || 'March'}</option>
                                                    <option value="4">{t('event.april') || 'April'}</option>
                                                    <option value="5">{t('event.may') || 'May'}</option>
                                                    <option value="6">{t('event.june') || 'June'}</option>
                                                    <option value="7">{t('event.july') || 'July'}</option>
                                                    <option value="8">{t('event.august') || 'August'}</option>
                                                    <option value="9">{t('event.september') || 'September'}</option>
                                                    <option value="10">{t('event.october') || 'October'}</option>
                                                    <option value="11">{t('event.november') || 'November'}</option>
                                                    <option value="12">{t('event.december') || 'December'}</option>
                                                </Form.Select>
                                            </Form.Group>
                                        </Col>
                                        <Col md={4}>
                                            <Form.Group className="mb-2">
                                                <Form.Label>{t('event.week_number') || 'Week'}</Form.Label>
                                                <Form.Select
                                                    name="yearly_week_number"
                                                    value={formData.yearly_week_number}
                                                    onChange={handleChange}
                                                >
                                                    <option value="0">{t('event.not_used') || 'Not used'}</option>
                                                    <option value="1">{t('event.first') || 'First'}</option>
                                                    <option value="2">{t('event.second') || 'Second'}</option>
                                                    <option value="3">{t('event.third') || 'Third'}</option>
                                                    <option value="4">{t('event.fourth') || 'Fourth'}</option>
                                                    <option value="5">{t('event.last') || 'Last'}</option>
                                                </Form.Select>
                                            </Form.Group>
                                        </Col>
                                        <Col md={4}>
                                            <Form.Group className="mb-2">
                                                <Form.Label>{t('event.day') || 'Day'}</Form.Label>
                                                <Form.Select
                                                    name="yearly_week_day"
                                                    value={formData.yearly_week_day}
                                                    onChange={handleChange}
                                                >
                                                    {weekDays.map(day => (
                                                        <option key={day.value} value={day.value}>{day.label}</option>
                                                    ))}
                                                </Form.Select>
                                            </Form.Group>
                                        </Col>
                                    </Row>

                                    <div className="text-center my-2 text-muted">{t('common.or') || 'OR'}</div>

                                    <Form.Group className="mb-2">
                                        <Form.Label>{t('event.day_of_year') || 'Day of the year'}</Form.Label>
                                        <Form.Control
                                            type="number"
                                            name="yearly_day_of_the_year"
                                            value={formData.yearly_day_of_the_year}
                                            onChange={handleChange}
                                            min="0"
                                            max="366"
                                        />
                                        <Form.Text className="text-muted">
                                            {t('event.day_of_year_help') || '0 = not used, 1-366'}
                                        </Form.Text>
                                    </Form.Group>
                                </div>
                            </Tab>
                        </Tabs>

                        {/* Recurrence Range */}
                        <div className="mt-3 p-2 border rounded">
                            <Form.Label><strong>{t('event.recurrence_range') || 'Recurrence Range'}</strong></Form.Label>
                            
                            <Form.Group className="mb-3">
                                <Form.Label>{t('event.recurrence_times') || 'Number of occurrences'}</Form.Label>
                                <Form.Control
                                    type="number"
                                    name="recurrence_times"
                                    value={formData.recurrence_times}
                                    onChange={handleChange}
                                    min="0"
                                />
                                <Form.Text className="text-muted">
                                    {t('event.recurrence_times_help') || '0 = infinite, N = repeat N times'}
                                </Form.Text>
                            </Form.Group>

                            <div className="text-center my-2 text-muted">{t('common.or') || 'OR'}</div>

                            <Form.Group className="mb-3">
                                <Form.Label>{t('event.recurrence_end') || 'End date'}</Form.Label>
                                <Form.Control
                                    type="datetime-local"
                                    name="recurrence_end_date"
                                    value={formData.recurrence_end_date}
                                    onChange={handleChange}
                                />
                            </Form.Group>
                        </div>
                    </div>
                </Collapse>
            </div>

            <ObjectLinkSelector
                value={formData.fk_obj_id}
                onChange={handleChange}
                classname="DBObject"
                fieldName="fk_obj_id"
                label={t('files.linked_object') || 'Linked Object'}
                required={false}
            />

            {error && (
                <Alert variant="danger" className="mb-3">
                    {error}
                </Alert>
            )}

            <div className="d-flex gap-2">
                <Button 
                    variant="primary" 
                    type="submit"
                    disabled={saving || isDeleted}
                >
                    {saving ? (
                        <>
                            <Spinner
                                as="span"
                                animation="border"
                                size="sm"
                                role="status"
                                aria-hidden="true"
                                className="me-2"
                            />
                            {t('common.saving')}
                        </>
                    ) : (
                        <>
                            <i className="bi bi-check-lg me-1"></i>
                            {t('common.save')}
                        </>
                    )}
                </Button>
                <Button 
                    variant="secondary" 
                    onClick={onCancel}
                    disabled={saving}
                >
                    <i className="bi bi-x-lg me-1"></i>
                    {t('common.cancel')}
                </Button>
                {onDelete && data.id && (
                    <Button 
                        variant="danger" 
                        onClick={onDelete}
                        disabled={saving}
                        className="ms-auto"
                    >
                        <i className="bi bi-trash me-1"></i>
                        {t('common.delete')}
                    </Button>
                )}
            </div>
        </Form>
    );
}


export function Events() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);

  const searchClassname = "DBEvent";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },
    { name: t("dbobjects.created") || "Created", attribute: "creation_date", type: "dateSelector" },
    { name: t("dbobjects.modified") || "Modified", attribute: "last_modify_date", type: "dateSelector" },
    { name: t("event.start") || "Start", attribute: "start_date", type: "dateSelector" },
    { name: t("event.end") || "End", attribute: "end_date", type: "dateSelector" },
  ];

  const orderBy = "start_date";

  const resultsColumns = [
    { name: t("event.start") || "Start Date", attribute: "start_date", type: "datetime", hideOnSmall: false },
    { name: t("event.end") || "End Date", attribute: "end_date", type: "datetime", hideOnSmall: false },
    // { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
    // { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    // { name: t("files.preview") || "File", attribute: "id", type: "imageView", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
    { name: t("dbobjects.created") || "Created", attribute: "creation_date", type: "dateSelector", hideOnSmall: true },
    { name: t("dbobjects.modified") || "Modified", attribute: "last_modify_date", type: "dateSelector", hideOnSmall: true },
  ]
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
    );
}
